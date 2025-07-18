package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/algorithm9/flash-deal/internal/app/apiserver"
	"github.com/algorithm9/flash-deal/pkg/logger"
	"github.com/algorithm9/flash-deal/pkg/meta"
)

var cfgFilePath string

func init() {
	flag.StringVar(&cfgFilePath, "config", "./conf.toml", "flash deal configuration file")
	flag.Parse()
}

//	@title			Flash Deal API
//	@version		1.0
//	@description	Flash Deal API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	https://lcslearn.top
//	@contact.email	algorithm9@163.com

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				format: "Bearer {token}"

//	@schemes	https
//	@host		lcslearn.top
//	@BasePath	/
//
// Swagger Declarative Comments Format DOC: https://github.com/swaggo/swag/blob/master/README.md
func main() {
	logger.Init()
	defer logger.Close()
	logger.L().Info().
		Str("version", meta.Version).
		Str("buildTime", meta.BuildTime).
		Str("buildTime", meta.BuildTime).
		Str("commit", meta.Commit).
		Msg("Starting flash-deal server...")

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv, cleanup, err := apiserver.InitServer(cfgFilePath)
	if err != nil {
		logger.L().Fatal().Msgf("failed to init app: %v", err)
		return
	}
	defer cleanup()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.L().Fatal().Msgf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	logger.L().Info().Msgf("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.L().Fatal().Msgf("Server forced to shutdown: ", err)
	}

	logger.L().Info().Msg("Flash-deal server exiting")
}
