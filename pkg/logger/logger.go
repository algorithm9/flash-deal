package logger

import (
	"io"
	"os"

	"github.com/phuslu/log"
)

var lg *log.Logger

// Init initializes the global logger.
func Init() {

	_ = os.MkdirAll("logs", 0755)

	env := os.Getenv("GO_ENV")
	var writer log.Writer
	fileWriter := &log.AsyncWriter{
		ChannelSize:   4096,
		DiscardOnFull: false,
		Writer: &log.FileWriter{
			Filename:   "logs/flashdeal.log",
			MaxSize:    50 * 1024 * 1024, // 50MB
			MaxBackups: 10,
		},
	}
	if env == "prod" {
		writer = fileWriter
	} else {
		writer = &log.IOWriter{Writer: os.Stderr}
	}

	log.DefaultLogger = log.Logger{
		Level:  log.InfoLevel,
		Caller: 1,
		Writer: writer,
	}

	log.Info().Str("env", os.Getenv("GO_ENV")).Msg("logger initialized")
}

func L() *log.Logger {
	return &log.DefaultLogger
}

func Close() {
	err := log.DefaultLogger.Writer.(io.Closer).Close()
	if err != nil {
		return
	}
}
