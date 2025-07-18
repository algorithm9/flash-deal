package config

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	LoadConfig,
	wire.Bind(new(ServerProvider), new(*Config)),
	wire.Bind(new(DBProvider), new(*Config)),
	wire.Bind(new(RedisProvider), new(*Config)),
	wire.Bind(new(LogProvider), new(*Config)),
	wire.Bind(new(MachineProvider), new(*Config)),
	wire.Bind(new(SMSProvider), new(*Config)),
	wire.Bind(new(JWTProvider), new(*Config)),
	wire.Bind(new(KafkaProvider), new(*Config)),
	wire.Bind(new(CacheProvider), new(*Config)),
	ProvideServer,
	ProvideDB,
	ProvideRedis,
	ProvideLog,
	ProvideMachine,
	ProvideSMS,
	ProvideJWT,
	ProvideKafka,
	ProvideCache,
)
