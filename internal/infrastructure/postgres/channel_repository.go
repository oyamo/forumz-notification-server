package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/oyamo/forumz-notification-server/internal/domain/channel"
)

type postgresChannelRepository struct {
	db *sql.DB
}

func (m postgresChannelRepository) FindChannels(ctx context.Context) ([]channel.Channel, error) {
	stmt, err := m.db.Prepare(`select c.id, c.channel_type_id, t.k, c.name, c.datetime_created, c.last_modified
			from channel c inner join channel_type t on c.channel_type_id = t.id`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var channels []channel.Channel
	for rows.Next() {
		var c channel.Channel
		err = rows.Scan(
			&c.ID,
			&c.ChannelTypeId,
			&c.ChannelTypeK,
			&c.Name,
			&c.DatetimeCreated,
			&c.LastModified)

		if err != nil {
			return nil, err
		}
		channels = append(channels, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (m postgresChannelRepository) InsertChannels(ctx context.Context, channel []channel.Channel) error {
	return errors.New("not implemented")
}

func NewPostgresChannelRepository(db *sql.DB) channel.Repository {
	return &postgresChannelRepository{
		db: db,
	}
}
