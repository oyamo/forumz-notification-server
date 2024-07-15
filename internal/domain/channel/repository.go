package channel

import "context"

type Repository interface {
	FindChannels(ctx context.Context) ([]Channel, error)
	InsertChannels(ctx context.Context, channel []Channel) error
}
