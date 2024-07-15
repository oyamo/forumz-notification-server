package main

import (
	"context"
	"github.com/oyamo/forumz-notification-server/internal/domain/channel"
	"github.com/oyamo/forumz-notification-server/internal/domain/notification"
	"github.com/oyamo/forumz-notification-server/internal/infrastructure/postgres"
	infraRedis "github.com/oyamo/forumz-notification-server/internal/infrastructure/redis"
	"github.com/oyamo/forumz-notification-server/internal/interfaces/kafka"
	"github.com/oyamo/forumz-notification-server/internal/interfaces/notifications"
	"github.com/oyamo/forumz-notification-server/internal/pkg"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

func main() {
	logConf, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	logger := logConf.Sugar()

	conf, err := pkg.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	conn, err := pkg.NewPostgresClient(conf.DatabaseDSN)
	defer conn.Close()
	if err != nil {
		logger.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: conf.RedisServer,
	})

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		logger.Fatalf("error dialing redis: %s\n", err)
	}

	notificationsRepo := postgres.NewPostgresQueueRepository(conn)
	channelsRepo := postgres.NewPostgresChannelRepository(conn)
	redisChannelRepo := infraRedis.NewRedisChannelRepository(redisClient)
	notificationsUsecase := notification.NewUsecase(logger, notificationsRepo)
	channelUsecase := channel.NewUsecase(logger, channelsRepo, redisChannelRepo)

	consumer := pkg.NewConsumer("notification-server", conf.KafkaConsumerServer)
	kafkaCRouter := kafka.NewConsumerRouter(consumer, logger, notificationsUsecase)
	kafkaCRouter.Consume()

	mailClient := pkg.NewMailClient(conf)
	sender := notifications.NewSender(logger, notificationsUsecase, channelUsecase, mailClient)
	go func() {
		sender.Start()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan
	logger.Errorw("received shutdown signal", "signal", sig)
}
