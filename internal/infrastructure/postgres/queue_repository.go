package postgres

import (
	"context"
	"database/sql"
	"github.com/oyamo/forumz-notification-server/internal/domain/notification"
)

type postgresNotificationRepository struct {
	db *sql.DB
}

func (p postgresNotificationRepository) Save(ctx context.Context, queue *notification.Notification) error {
	stmt, err := p.db.Prepare(`insert into notification(id, email_address, content, priority, type) 
			values ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(
		ctx,
		queue.ID,
		queue.EmailAddress,
		queue.Content,
		queue.Priority,
		queue.Type,
	)

	if err != nil {
		return err
	}
	return nil
}

func (p postgresNotificationRepository) FindPending(ctx context.Context) ([]notification.Notification, error) {
	stmt, err := p.db.Prepare(`select q.id, q.email_address, q.content, q.priority, q.retry_count, q.sent, q.type, q.datetime_created, q.last_modified 
			from notification q where q.sent = false order by priority , retry_count , datetime_created  , last_modified`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var queues []notification.Notification
	if row.Next() {
		q := notification.Notification{}
		err = row.Scan(
			&q.ID,
			&q.EmailAddress,
			&q.Content,
			&q.Priority,
			&q.AttemptCount,
			&q.Sent,
			&q.Type,
			&q.DatetimeCreated,
			&q.LastModified,
		)

		if err != nil {
			return nil, err
		}

		queues = append(queues, q)
	}

	return queues, nil
}

func (p postgresNotificationRepository) Update(ctx context.Context, queue *notification.Notification) error {
	stmt, err := p.db.Prepare(`update notification set sent = $2, retry_count = $3  where id = $1`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.ExecContext(
		ctx,
		queue.ID,
		queue.Sent,
		queue.AttemptCount,
	)
	if err != nil {
		return err
	}

	return nil
}

func NewPostgresQueueRepository(db *sql.DB) notification.Repository {
	return &postgresNotificationRepository{
		db: db,
	}
}
