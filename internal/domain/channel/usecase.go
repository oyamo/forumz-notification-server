package channel

import (
	"context"
	"go.uber.org/zap"
)

type Usecase struct {
	logger           *zap.SugaredLogger
	redisChannelRepo Repository
	postgresChanRepo Repository
}

func (u *Usecase) FindChannels(ctx context.Context) ([]Channel, error) {
	channels, err := u.redisChannelRepo.FindChannels(ctx)
	if err == nil {
		return channels, nil
	}

	channels, err = u.postgresChanRepo.FindChannels(ctx)
	if err == nil {
		err = u.redisChannelRepo.InsertChannels(ctx, channels)
		if err == nil {
			return channels, nil
		}
		return nil, err
	}

	return nil, err
}

func NewUsecase(logger *zap.SugaredLogger, redisChannelRepo, postgresChanRepo Repository) *Usecase {
	return &Usecase{
		logger:           logger,
		redisChannelRepo: redisChannelRepo,
		postgresChanRepo: postgresChanRepo,
	}
}
