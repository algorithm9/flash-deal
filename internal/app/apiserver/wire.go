//go:build wireinject
// +build wireinject

package apiserver

import (
	"github.com/google/wire"

	"github.com/algorithm9/flash-deal/internal/config"
	"github.com/algorithm9/flash-deal/internal/module/product"
	"github.com/algorithm9/flash-deal/internal/module/seckill"
	"github.com/algorithm9/flash-deal/internal/module/user"
	"github.com/algorithm9/flash-deal/internal/shared"
	"github.com/algorithm9/flash-deal/pkg/cache"
)

func InitServer(cfgPath string) (*Server, func(), error) {
	panic(
		wire.Build(
			config.ProviderSet,
			shared.ProviderSet,
			user.ProviderSet,
			product.ProviderSet,
			ProviderSet,
			seckill.ProviderSet,
			cache.Provider,
		),
	)

	return &Server{}, nil, nil
}
