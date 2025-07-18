package idgen

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSnowflakeIDGen,
	wire.Bind(new(IDGenerator), new(*SnowflakeIDGen)),
)
