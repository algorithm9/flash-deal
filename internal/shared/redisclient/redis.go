package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/algorithm9/flash-deal/internal/model"
	"github.com/algorithm9/flash-deal/pkg/errorx"
	"github.com/algorithm9/flash-deal/pkg/logger"
)

type Client struct {
	Client *redis.Client
}

func NewClient(conf *model.RedisConfig) (*Client, func()) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.URL,
		Password: conf.Passwd,
		DB:       conf.DB,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.L().Fatal().Err(err).Msg("failed to init redisclient")
	}
	logger.L().Info().Msg("Redis client created successfully.")

	cleanup := func() {
		logger.L().Info().Msg("closing redis client")
		if err := redisClient.Close(); err != nil {
			logger.L().Err(err).Msgf("redis close error")
		}
	}
	return &Client{Client: redisClient}, cleanup
}

func IsNil(err error) bool {
	return errorx.Is(err, redis.Nil)
}
