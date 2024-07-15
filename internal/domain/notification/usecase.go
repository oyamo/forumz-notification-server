package notification

import (
	"context"
	"go.uber.org/zap"
)

type Usecase struct {
	logger          *zap.SugaredLogger
	queueRepository Repository
}

func (u Usecase) FindPending(ctx context.Context) ([]Notification, error) {
	queue, err := u.queueRepository.FindPending(ctx)
	if err != nil {
		return nil, err
	}
	return queue, nil
}

func (u Usecase) Update(ctx context.Context, queue *Notification) error {
	err := u.queueRepository.Update(ctx, queue)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) Save(ctx context.Context, queue *Notification) error {
	err := u.queueRepository.Save(ctx, queue)
	if err != nil {
		return err
	}
	return nil
}

func NewUsecase(logger *zap.SugaredLogger, queueRepository Repository) *Usecase {
	return &Usecase{
		logger:          logger,
		queueRepository: queueRepository,
	}
}
