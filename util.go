package main

import (
	"github.com/optimizely/go-sdk/pkg/client"
	"time"
)

func initOptimizely() *client.OptimizelyClient {
	optimizelyFactory := &client.OptimizelyFactory{
		SDKKey: "PEEY3bqav2mZXVJWkQWsm",
	}

	datafilePollingInterval := 2 * time.Second
	eventBatchSize := 20
	eventQueueSize := 1500
	eventFlushInterval := 10 * time.Second

	// Instantiate a client with custom configuration
	oC, _ := optimizelyFactory.Client(
		client.WithPollingConfigManager(datafilePollingInterval, nil),
		client.WithBatchEventProcessor(
			eventBatchSize,
			eventQueueSize,
			eventFlushInterval,
		),
	)
	return oC
}

