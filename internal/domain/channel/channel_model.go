package channel

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type ChannelName string

const (
	Mailtrap = "Mailtrap"
	Sendgrid = "SendGrid"
)

type Channel struct {
	ID              uuid.UUID   `json:"id"`
	ChannelTypeId   string      `json:"channelTypeId"`
	ChannelTypeK    string      `json:"channelTypeK"`
	Name            ChannelName `json:"name"`
	DatetimeCreated time.Time   `json:"datetimeCreated"`
	LastModified    time.Time   `json:"lastModified"`
}

func (c *Channel) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}
