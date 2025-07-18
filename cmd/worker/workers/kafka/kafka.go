package kafka

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	_ "github.com/go-sql-driver/mysql"

	db "github.com/algorithm9/flash-deal/cmd/worker/workers/db/models"
	"github.com/algorithm9/flash-deal/internal/model"
	"github.com/algorithm9/flash-deal/internal/shared/redisclient"
	"github.com/algorithm9/flash-deal/pkg/logger"
)

type Consumer struct {
	consumer    *kafka.Consumer
	db          *sql.DB
	queries     *db.Queries
	rdb         *redisclient.Client
	workerCount int
	pollTimeout time.Duration
	queueSize   int
}

func NewConsumer(kafkaConfig *model.Kafka, dbConfig *model.DatabaseConfig, rdbCfg *model.RedisConfig) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Server,
		"group.id":          kafkaConfig.GroupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	if err := c.SubscribeTopics([]string{kafkaConfig.Topic}, nil); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		dbConfig.UserName, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName, "UTC",
	)
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.L().Fatal().Err(err).Msg("failed to open database")
	}
	if err := database.Ping(); err != nil {
		database.Close()
		c.Close()
		logger.L().Fatal().Err(err).Msg("failed to connect to database")
	}

	database.SetMaxOpenConns(dbConfig.MaxOpenConns)
	database.SetMaxIdleConns(dbConfig.MaxIdleConns)

	q := db.New(database)

	rdb, _ := redisclient.NewClient(rdbCfg)

	return &Consumer{
		consumer:    c,
		db:          database,
		queries:     q,
		rdb:         rdb,
		workerCount: kafkaConfig.WorkerCount,
		pollTimeout: time.Duration(kafkaConfig.PollTimeout) * time.Millisecond,
		queueSize:   kafkaConfig.QueueSize,
	}, nil
}
