package kafka

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	seckilldto "github.com/algorithm9/flash-deal/internal/module/seckill/dto"
	"github.com/algorithm9/flash-deal/pkg/logger"
)

func (kc *Consumer) Run(ctx context.Context) {

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	msgChan := make(chan *kafka.Message, kc.queueSize)
	defer close(msgChan)

	go func() {
		for runCtx.Err() == nil {
			msg, err := kc.consumer.ReadMessage(kc.pollTimeout)

			if err != nil {
				if err.(kafka.Error).Code() != kafka.ErrTimedOut {
					logger.L().Error().Err(err).Msg("Kafka message read failed")
				}
				continue
			}

			msgChan <- msg
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < kc.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range msgChan {
				var req seckilldto.SeckillRequest
				if err := json.Unmarshal(msg.Value, &req); err != nil {
					continue
				}
				err := kc.handleOrder(ctx, req)
				if err != nil {
					logger.L().Err(err).Msgf("Failed to handle order")
					// todo retry and dlq
					continue
				}
			}
		}()
	}

	<-ctx.Done()
	cancel()
	wg.Wait()

	kc.consumer.Close()
	kc.db.Close()
	kc.rdb.Client.Close()
}
