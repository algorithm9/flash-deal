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
	DB    model.DatabaseConfig `toml:"database"`
	Redis model.RedisConfig    `toml:"redis"`
	Kafka model.Kafka          `toml:"kafka"`
}
