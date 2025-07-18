package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/algorithm9/flash-deal/internal/model"
)

type Options struct {
	addr   string
	logCfg *model.LogConfig
}

func NewOptions(cfg *model.ServerConfig, logCfg *model.LogConfig) *Options {
	return &Options{
		addr:   cfg.Addr,
		logCfg: logCfg,
	}
}

type Server struct {
	*http.Server
}

func New(opts *Options, r *gin.Engine) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    opts.addr,
			Handler: r,
		},
	}
}

var ProviderSet = wire.NewSet(New, NewOptions, NewGinEngine)
