package notifications

import (
	"context"
	"errors"
	"github.com/oyamo/forumz-notification-server/internal/domain/channel"
	"github.com/oyamo/forumz-notification-server/internal/domain/notification"
	"github.com/oyamo/forumz-notification-server/internal/pkg"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

type Sender struct {
	notificationsUseCase *notification.Usecase
	channelUseCase       *channel.Usecase
	ticker               *time.Ticker
	channels             []channel.Channel
	channelAtomicIdx     atomic.Int32
	mailSender           *pkg.MailClient
	logger               *zap.SugaredLogger
}

func (s *Sender) Start() {

	dispatcher := NewDispatcher(5, s.sendNotification, s.updateNotification)

	go func() {
		dispatcher.Start()
	}()

	for {
		select {
		case <-s.ticker.C:
			ctx := context.Background()
			var emailChannels []channel.Channel
			channels, err := s.channelUseCase.FindChannels(ctx)
			if err != nil {
				s.logger.Error(err)
				return
			}

			for _, notificationChannel := range channels {
				if notificationChannel.ChannelTypeK == "Mail" {
					emailChannels = append(emailChannels, notificationChannel)
				}
			}

			if len(emailChannels) <= 0 {
				continue
			}

			s.channels = emailChannels

			notifications, err := s.notificationsUseCase.FindPending(ctx)
			if err != nil {
				s.logger.Error(err)
				continue
			}

			for _, n := range notifications {
				dispatcher.AddToWorkerPool(&n)
			}
		}
	}

}

func (s *Sender) sendNotification(notification *notification.Notification) error {
	if len(s.channels) <= 0 {
		err := errors.New("no channels found")
		s.logger.Error(err)
		return err
	}

	index := s.channelAtomicIdx.Add(1)
	notificationChannel := s.channels[index%int32(len(s.channels))]

	err := s.mailSender.SendMail(notification, notificationChannel.Name)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}

func (s *Sender) updateNotification(notification *notification.Notification) error {
	ctx := context.Background()
	err := s.notificationsUseCase.Update(ctx, notification)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func NewSender(logger *zap.SugaredLogger, notificationsUseCase *notification.Usecase, channelUseCase *channel.Usecase, mailSender *pkg.MailClient) *Sender {
	return &Sender{
		notificationsUseCase: notificationsUseCase,
		channelUseCase:       channelUseCase,
		ticker:               time.NewTicker(time.Second * 5),
		logger:               logger,
		mailSender:           mailSender,
	}
}
