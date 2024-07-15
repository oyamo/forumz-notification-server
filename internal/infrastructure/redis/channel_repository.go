package redis

import (
	"context"
	"encoding/json"
	"github.com/oyamo/forumz-notification-server/internal/domain/channel"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	CacheDuration = time.Minute * 2
)

type redisChannelRepository struct {
	client *redis.Client
}

func (r redisChannelRepository) InsertChannels(ctx context.Context, channels []channel.Channel) error {
	// marshalInto Bytes
	b, err := json.Marshal(&channels)
	if err != nil {
		return err
	}

	_, err = r.client.Set(ctx, "channels", b, CacheDuration).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r redisChannelRepository) FindChannels(ctx context.Context) ([]channel.Channel, error) {
	res, err := r.client.Get(ctx, "channels").Result()
	if err != nil {
		return nil, err
	}
	var channels []channel.Channel
	err = json.Unmarshal([]byte(res), &channels)
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func NewRedisChannelRepository(client *redis.Client) channel.Repository {
	return &redisChannelRepository{
		client: client,
	}
}
