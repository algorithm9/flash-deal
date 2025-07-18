package product

import (
	"github.com/google/wire"

	producthttp "github.com/algorithm9/flash-deal/internal/module/product/delivery/http"
	"github.com/algorithm9/flash-deal/internal/module/product/repository"
	"github.com/algorithm9/flash-deal/internal/module/product/service"
)

var ProviderSet = wire.NewSet(
	repository.Provider,
	service.Provider,
	producthttp.Provider,
)
