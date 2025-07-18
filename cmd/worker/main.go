package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/algorithm9/flash-deal/cmd/worker/app"
	"github.com/algorithm9/flash-deal/cmd/worker/config"
	"github.com/algorithm9/flash-deal/pkg/logger"
	"github.com/algorithm9/flash-deal/pkg/meta"
)

var (
	cfgFilePath string
)

func init() {
	flag.StringVar(&cfgFilePath, "config", "./conf.toml", "flash deal configuration file")
	flag.Parse()
}

func main() {
	logger.Init()
	defer logger.Close()
	logger.L().Info().
		Str("version", meta.Version).
		Str("buildTime", meta.BuildTime).
		Str("buildTime", meta.BuildTime).
		Str("commit", meta.Commit).
		Msg("Starting flash-deal worker...")

	cfg := config.LoadConfig(cfgFilePath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// catching exit signals
	go handleSignal(cancel)

	// start all workers
	if err := app.Run(ctx, cfg); err != nil {
		logger.L().Fatal().Err(err).Msg("Worker run error:")
	}
	logger.L().Info().Msg("Worker service exited gracefully.")
}

func handleSignal(cancelFunc context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	logger.L().Info().Msg("Shutdown signal received.")
	cancelFunc()
}
