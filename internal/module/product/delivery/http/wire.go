package producthttp

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewGoodsHandler,
	NewGoodsRouter,
)
