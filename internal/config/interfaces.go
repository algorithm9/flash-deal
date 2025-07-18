package config

import "github.com/algorithm9/flash-deal/internal/model"

type ServerProvider interface {
	ProvideServer() *model.ServerConfig
}

func ProvideServer(s ServerProvider) *model.ServerConfig {
	return s.ProvideServer()
}

type DBProvider interface {
	ProvideDB() *model.DatabaseConfig
}

func ProvideDB(p DBProvider) *model.DatabaseConfig {
	return p.ProvideDB()
}

type RedisProvider interface {
	ProvideRedis() *model.RedisConfig
}

func ProvideRedis(p RedisProvider) *model.RedisConfig {
	return p.ProvideRedis()
}

type LogProvider interface {
	ProvideLog() *model.LogConfig
}

func ProvideLog(p LogProvider) *model.LogConfig {
	return p.ProvideLog()
}

type MachineProvider interface {
	ProvideMachine() *model.Machine
}

func ProvideMachine(p MachineProvider) *model.Machine {
	return p.ProvideMachine()
}

type SMSProvider interface {
	ProvideSMS() *model.SMS
}

func ProvideSMS(p SMSProvider) *model.SMS {
	return p.ProvideSMS()
}

type JWTProvider interface {
	ProvideJWT() *model.JWT
}

func ProvideJWT(p JWTProvider) *model.JWT {
	return p.ProvideJWT()
}

type KafkaProvider interface {
	ProvideKafka() *model.Kafka
}

func ProvideKafka(p KafkaProvider) *model.Kafka {
	return p.ProvideKafka()
}

type CacheProvider interface {
	ProvideCache() *model.Cache
}

func ProvideCache(p CacheProvider) *model.Cache {
	return p.ProvideCache()
}
