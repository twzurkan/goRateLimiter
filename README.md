# goRateLimiter

The goRateLimiter uses a ratelimiter.go which uses the limit and duration to increment redis key on ip.
You can simply clone this repo and run docker-compose up in this directory.
The, go the http://localhost:5000/get to see the get method etc.

It also uses a feature and experiment from Optimizely.  The current SDK key is for my personal account.  [Optimizely](https://optimizely.com) provides feature flags and 1 experiment with a feature flag for free.

Right now, it uses these attributes "ip", "time", "limit".  As an example you can do the following to test:
As long as you have a docker client.  You can do the following:
1. docker-compose up
2. curl -X POST http://localhost:5000/post ; curl -H "X-location:CA" -H "X-time:7:00 pm" -X POST http://localhost:5000/post
curl http://localhost:5000/get ; curl -H "X-location:CA" -H "X-time:7:00 pm" http://localhost:5000/get
The above includes two custom headers which must both exist in order for the ip to be allowed 20 request per 2 seconds.  The default is 10 per 1 second.

