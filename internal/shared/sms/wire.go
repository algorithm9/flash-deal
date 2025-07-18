package sms

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewSMSClient,
)
