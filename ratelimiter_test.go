package main

import (
	"testing"
)

// TestNewRateLimiter makes sure that the redis instance is created.
func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter()
	if rl.rdb == nil {
		t.Error("rdb is nil. start redis first.")
	}
}

// TestRateLimiter_Throttle tests 1000 messages in 1 second. and should throttle.
func TestRateLimiter_Throttle(t *testing.T) {
   rl := NewRateLimiter()
   throttleCount := 0
   for i := 0; i < 1000; i++ {
   		if rl.Throttle( "TestRateLimiter_Throttle", 1000, 1) {
   			throttleCount++
		}
   }

   if throttleCount > 0 {
   		t.Failed()
   }

}

// TestRateLimiter_Throttle tests 1000 messages in 1 second. and should throttle.
func TestRateLimiter_Throttle_Fail(t *testing.T) {
	rl := NewRateLimiter()
	throttleCount := 0
	for i := 0; i < 1000; i++ {
		if rl.Throttle( "TestRateLimiter_Throttle_Fail", 900, 1) {
			throttleCount++
		}
	}

	if throttleCount != 100 {
		t.Failed()
	}

}
