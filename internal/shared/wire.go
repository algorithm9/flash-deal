package shared

import (
	"github.com/google/wire"

	"github.com/algorithm9/flash-deal/internal/shared/entx"
	"github.com/algorithm9/flash-deal/internal/shared/idgen"
	"github.com/algorithm9/flash-deal/internal/shared/middleware"
	"github.com/algorithm9/flash-deal/internal/shared/redisclient"
	"github.com/algorithm9/flash-deal/internal/shared/sms"
)

var ProviderSet = wire.NewSet(
	idgen.ProviderSet,
	sms.Provider,
	redisclient.Provider,
	entx.NewEntClient,
	middleware.ProviderSet,
)
