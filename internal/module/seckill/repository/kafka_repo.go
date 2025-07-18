package repository

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/algorithm9/flash-deal/internal/model"
	seckilldto "github.com/algorithm9/flash-deal/internal/module/seckill/dto"
	"github.com/algorithm9/flash-deal/pkg/logger"
)

type kafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(k *model.Kafka) (QueueProducer, func(), error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": k.Server})
	if err != nil {
		panic(err)
	}

	// 异步 delivery handler
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.L().Info().Msgf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					logger.L().Info().Msgf("Delivered to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// 清理函数
	cleanup := func() {
		p.Close()
		logger.L().Info().Msg("closing kafka producer client")
	}

	return &kafkaProducer{
		producer: p,
		topic:    k.Topic,
	}, cleanup, nil
}

func (k *kafkaProducer) Send(ctx context.Context, msg seckilldto.SeckillRequest) error {
	return nil
	//value, err := json.Marshal(msg)
	//if err != nil {
	//	return err
	//}
	//return k.producer.Produce(&kafka.Message{
	//	TopicPartition: kafka.TopicPartition{Topic: &k.topic, Partition: kafka.PartitionAny},
	//	Key:            []byte(fmt.Sprintf("%d", msg.UserID)),
	//	Value:          value,
	//}, nil)
}
