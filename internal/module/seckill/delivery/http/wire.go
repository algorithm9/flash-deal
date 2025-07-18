package seckillhttp

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewSeckillRouter,
	NewSeckillHandler,
)
