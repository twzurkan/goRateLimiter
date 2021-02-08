package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

var ctx = context.Background()

type RateLimiter struct {
	rdb *redis.Client
}

func (l* RateLimiter) attach() {
	l.rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func NewRateLimiter() *RateLimiter {
	r := &RateLimiter{}
	r.attach()

	return r
}

func (l *RateLimiter) incrementAndSetExpire(ip string, duration time.Duration) int64 {
	var incr *redis.IntCmd
	_, err := l.rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		incr = pipe.Incr(ctx, ip)
		pipe.Expire(ctx, ip, duration*time.Second)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// The value is available only after pipeline is executed.
	return incr.Val()
}

// Limit limit the number of time the ip can hit the endpoint.
func (l *RateLimiter) Throttle(ip string, limit int, duration int) bool {

	_, err := l.rdb.Get(ctx, ip).Result()
	var val int64 = 0
	if err == redis.Nil {
		val = l.incrementAndSetExpire(ip, time.Duration(duration))
		fmt.Printf("key does not exist, %v\n", ip)
	} else if err != nil {
		panic(err)
	} else {
		v := l.rdb.Incr(ctx, ip)
		fmt.Printf("here is %d\n", v)
		val = v.Val()
	}

	if val > int64(limit) {
		return true
	}

	return false
}
