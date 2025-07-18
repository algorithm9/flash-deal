package seckill

import (
	"github.com/google/wire"

	seckillhttp "github.com/algorithm9/flash-deal/internal/module/seckill/delivery/http"
	"github.com/algorithm9/flash-deal/internal/module/seckill/repository"
	"github.com/algorithm9/flash-deal/internal/module/seckill/service"
)

var ProviderSet = wire.NewSet(
	seckillhttp.Provider,
	service.Provider,
	repository.Provider,
)
