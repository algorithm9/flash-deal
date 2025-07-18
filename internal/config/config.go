package config

import (
	"github.com/BurntSushi/toml"

	"github.com/algorithm9/flash-deal/internal/model"
)

func LoadConfig(path string) *Config {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}

type Config struct {
	Server   model.ServerConfig   `toml:"server"`
	Database model.DatabaseConfig `toml:"database"`
	Redis    model.RedisConfig    `toml:"redisclient"`
	Log      model.LogConfig      `toml:"log"`
	Machine  model.Machine        `toml:"machine"`
	SMS      model.SMS            `toml:"sms"`
	JWT      model.JWT            `toml:"jwt"`
	Kafka    model.Kafka          `toml:"kafka"`
	Cache    model.Cache          `toml:"cache"`
}

func (cfg *Config) ProvideServer() *model.ServerConfig {
	return &cfg.Server
}

func (cfg *Config) ProvideDB() *model.DatabaseConfig {
	return &cfg.Database
}

func (cfg *Config) ProvideRedis() *model.RedisConfig {
	return &cfg.Redis
}

func (cfg *Config) ProvideLog() *model.LogConfig {
	return &cfg.Log
}

func (cfg *Config) ProvideMachine() *model.Machine {
	return &cfg.Machine
}

func (cfg *Config) ProvideSMS() *model.SMS {
	return &cfg.SMS
}

func (cfg *Config) ProvideJWT() *model.JWT {
	return &cfg.JWT
}

func (cfg *Config) ProvideKafka() *model.Kafka {
	return &cfg.Kafka
}

func (cfg *Config) ProvideCache() *model.Cache {
	return &cfg.Cache
}
