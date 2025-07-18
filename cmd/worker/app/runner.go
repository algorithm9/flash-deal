package app

import (
	"context"
	"sync"

	"github.com/algorithm9/flash-deal/cmd/worker/config"
)

// Run start all workers
func Run(ctx context.Context, cfg *config.Config) error {
	workers := RegisterWorkers(cfg)

	var wg sync.WaitGroup
	for _, w := range workers {
		wg.Add(1)
		go func(w Worker) {
			defer wg.Done()
			w.Run(ctx)
		}(w)
	}

	<-ctx.Done()
	wg.Wait()
	return nil
}
