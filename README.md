# goRateLimiter

The goRateLimiter uses a ratelimiter.go which uses the limit and duration to set the redis instance count.
You can simply clone this repo and run docker-compose up in this directory.
The, go the http://localhost:5000/get to see the get method etc.

It also uses a feature and experiment from Optimizely.  The current SDK key is for my personal account.  Optimizely provides feature flags and 1 experiment with a feature flag for free.
