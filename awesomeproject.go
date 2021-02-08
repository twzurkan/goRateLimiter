package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/optimizely/go-sdk/pkg/entities"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var rateLimiter = NewRateLimiter()
var optimizelyClient = initOptimizely()

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/get", getHandler)
	r.HandleFunc("/post", postHandler)
	r.Use(rateLimitMiddleware)
	r.Use(loggingMiddleware)

	srv := &http.Server{
		Addr:         "0.0.0.0:5000",
		// Good practice to set timeouts
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func getClientIp(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIp(r)

		attributes := map[string]interface{}{
			"ip":      ip,
		}

		locations := r.Header["X-Location"]
		times := r.Header["X-Time"]
		if locations != nil {
			attributes["location"] = locations[0]
		}
		if times != nil {
			attributes["time"] = times[0]
		}
		
		user := entities.UserContext{
			ID:         "userId",
			Attributes: attributes,
		}
		duration, _ := optimizelyClient.GetFeatureVariableInteger("rate_limit", "time", user)
		limit, _ := optimizelyClient.GetFeatureVariableInteger("rate_limit", "limit", user)
		fmt.Printf("ip: %v, duration: %v, limit: %v\n", ip, duration, limit)
		if rateLimiter.Throttle(ip, limit, duration) == true {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, http.StatusText(429), 429)
			err := optimizelyClient.Track("over_limit", user, map[string]interface{}{})
			if err != nil {
				fmt.Printf("%v", err)
			}
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	_, err := fmt.Fprint(w, "You called get")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	_, err := fmt.Fprint(w, "You posted")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}