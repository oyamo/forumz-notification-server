package notification

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, queue *Notification) error
	FindPending(ctx context.Context) ([]Notification, error)
	Update(ctx context.Context, queue *Notification) error
}
