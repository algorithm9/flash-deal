package middleware

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAuthJWTMiddleware,
)
