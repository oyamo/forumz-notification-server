package notification

import (
	"github.com/google/uuid"
	"time"
)

type NotificationType string

const (
	NotificationTypeConnection NotificationType = "Connection"
	NotificationTypePost                        = "Post"
)

type Notification struct {
	ID              uuid.UUID
	EmailAddress    string
	Content         string
	Type            NotificationType
	Priority        int
	AttemptCount    int
	Sent            bool
	DatetimeCreated time.Time
	LastModified    time.Time
}
