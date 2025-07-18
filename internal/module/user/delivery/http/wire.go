package userhttp

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewAuthHandler,
	NewUserRouter,
)
