package app

import (
	"context"

	"github.com/algorithm9/flash-deal/cmd/worker/config"
	"github.com/algorithm9/flash-deal/cmd/worker/workers/kafka"
)

type Worker interface {
	Run(ctx context.Context)
}

// RegisterWorkers return all worker instances
func RegisterWorkers(cfg *config.Config) []Worker {
	kafkaWorker, err := kafka.NewConsumer(&cfg.Kafka, &cfg.DB, &cfg.Redis)
	if err != nil {
		panic(err)
	}
	return []Worker{
		kafkaWorker,
	}
}
