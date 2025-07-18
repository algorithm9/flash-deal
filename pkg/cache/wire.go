package cache

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewCache,
)
