package pkg

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseDSN         string
	RedisServer         string
	KafkaConsumerServer string
	KafkaProducerServer string
	SenderName          string
	SenderEmail         string
	SendGridAPIKey      string
	MailtrapAPIKey      string
}

const (
	envDatabaseDSN    = "NOTIFICATIONS_SERVICE_DATABASE_DSN"
	envRedisServer    = "NOTIFICATIONS_SERVICE_REDIS_SERVER"
	envKafkaConsumer  = "NOTIFICATIONS_SERVICE_KAFKA_CONSUMER"
	envKafkaProducer  = "NOTIFICATIONS_SERVICE_KAFKA_PRODUCER"
	envSenderName     = "NOTIFICATIONS_SERVICE_SENDER_NAME"
	envSenderEmail    = "NOTIFICATIONS_SERVICE_SENDER_EMAIL"
	envSendGridAPIKey = "NOTIFICATIONS_SERVICE_SENDGRID_API_KEY"
	envMailtrapAPIKey = "NOTIFICATIONS_SERVICE_MAILTRAP_API_KEY"
)

func EnvNotSetError(env string) error {
	return fmt.Errorf("%s environment variable not set", env)
}

func NewConfig() (*Config, error) {
	dbDSN, ok := os.LookupEnv(envDatabaseDSN)
	if !ok {
		return nil, EnvNotSetError(envDatabaseDSN)
	}

	kafkaProducer, ok := os.LookupEnv(envKafkaProducer)
	if !ok {
		return nil, EnvNotSetError(envKafkaProducer)
	}

	kafkaConsumer, ok := os.LookupEnv(envKafkaConsumer)
	if !ok {
		return nil, EnvNotSetError(envKafkaConsumer)
	}

	redisServer, ok := os.LookupEnv(envRedisServer)
	if !ok {
		return nil, EnvNotSetError(envRedisServer)
	}

	senderName, ok := os.LookupEnv(envSenderName)
	if !ok {
		return nil, EnvNotSetError(envSenderName)
	}

	senderEmail, ok := os.LookupEnv(envSenderEmail)
	if !ok {
		return nil, EnvNotSetError(envSenderEmail)
	}

	sendGridAPIKey, ok := os.LookupEnv(envSendGridAPIKey)
	if !ok {
		return nil, EnvNotSetError(envSendGridAPIKey)
	}

	mailtrapAPIKey, ok := os.LookupEnv(envMailtrapAPIKey)
	if !ok {
		return nil, EnvNotSetError(envMailtrapAPIKey)
	}

	return &Config{
		DatabaseDSN:         dbDSN,
		RedisServer:         redisServer,
		KafkaConsumerServer: kafkaConsumer,
		KafkaProducerServer: kafkaProducer,
		SenderName:          senderName,
		SenderEmail:         senderEmail,
		SendGridAPIKey:      sendGridAPIKey,
		MailtrapAPIKey:      mailtrapAPIKey,
	}, nil
}
