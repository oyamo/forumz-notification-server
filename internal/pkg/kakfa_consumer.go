package pkg

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"log"
)

type HandleFunc func([]byte)

type KakfaConsumer struct {
	groupId        string
	consumerServer string
	logger         *zap.SugaredLogger
}

func (k *KakfaConsumer) ConsumeAndHandle(topic string, handler HandleFunc) {
	go func() {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{k.consumerServer},
			GroupID:  k.groupId,
			Topic:    topic,
			MaxBytes: 10e6, // 10MB
		})
		defer reader.Close()

		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			handler(msg.Value)

			if err := reader.CommitMessages(context.Background(), msg); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}()
}

func NewConsumer(groupId, consumer string) *KakfaConsumer {
	return &KakfaConsumer{
		groupId:        groupId,
		consumerServer: consumer,
	}
}
