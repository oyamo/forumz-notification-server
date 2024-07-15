package notification

import "time"

type NotificationDTO struct {
	Type            NotificationType  `json:"type"`
	Recipient       string            `json:"recipient"`
	AdditionalInfo  map[string]string `json:"additionalInfo,omitempty"`
	DatetimeCreated time.Time         `json:"datetimeCreated,omitempty"`
}
