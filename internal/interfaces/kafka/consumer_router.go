package kafka

import (
	"github.com/oyamo/forumz-notification-server/internal/domain/notification"
	"github.com/oyamo/forumz-notification-server/internal/interfaces/kafka/handlers"
	"github.com/oyamo/forumz-notification-server/internal/pkg"
	"go.uber.org/zap"
)

type ConsumerRouter struct {
	consumer        *pkg.KakfaConsumer
	logger          *zap.SugaredLogger
	notificationsUC *notification.Usecase
}

func (r *ConsumerRouter) Consume() {
	notificationsHandler := handlers.NewNotificationHandler(r.logger, r.notificationsUC)
	r.consumer.ConsumeAndHandle("Put-Notification-v1", notificationsHandler.HandleNotification)
}

func NewConsumerRouter(consumer *pkg.KakfaConsumer, logger *zap.SugaredLogger, notificationsUC *notification.Usecase) *ConsumerRouter {
	return &ConsumerRouter{
		consumer:        consumer,
		logger:          logger,
		notificationsUC: notificationsUC,
	}
}
