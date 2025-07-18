package user

import (
	"github.com/google/wire"

	userhttp "github.com/algorithm9/flash-deal/internal/module/user/delivery/http"
	"github.com/algorithm9/flash-deal/internal/module/user/repository"
	"github.com/algorithm9/flash-deal/internal/module/user/service"
)

var ProviderSet = wire.NewSet(
	userhttp.Provider,
	repository.ProviderSet,
	service.ProviderSet,
	//repository.NewUserRepository,
)
