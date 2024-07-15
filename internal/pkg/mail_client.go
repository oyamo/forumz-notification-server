package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/oyamo/forumz-notification-server/internal/domain/channel"
	"github.com/oyamo/forumz-notification-server/internal/domain/notification"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
	"time"
)

const (
	sendGridBaseURL      = "https://send.api.mailtrap.io"
	sendGridSendEndpoint = "/api/send"
)

type MailClient struct {
	sendgridClient *sendgrid.Client
	conf           *Config
	httpClient     *http.Client
}

type Person struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type MailRequest struct {
	To      []Person `json:"to"`
	From    Person   `json:"from"`
	Headers struct {
		XMessageSource string `json:"X-Message-Source"`
	} `json:"headers"`
	Subject  string `json:"subject"`
	Text     string `json:"text"`
	HTML     string `json:"html"`
	Category string `json:"category"`
}

func (c *MailClient) SendMail(payload *notification.Notification, channelName channel.ChannelName) error {
	var subject string
	switch payload.Type {
	case notification.NotificationTypeConnection:
		subject = "You have a new Connection"
	case notification.NotificationTypePost:
		subject = "Check out these posts"
	}
	switch channelName {
	case channel.Mailtrap:
		return c.submitMailtrap(subject, payload.EmailAddress, payload.Content, "text/html")
	case channel.Sendgrid:
		return c.submitSendgrid(subject, payload.EmailAddress, payload.Content, "text/html")
	default:
		return fmt.Errorf("invalid channel name %s", channelName)
	}
}

func (c *MailClient) submitMailtrap(subject, to string, content, contentType string) error {
	var mailRequest MailRequest
	mailRequest.Subject = subject
	mailRequest.HTML = content
	mailRequest.Headers.XMessageSource = ""

	mailRequest.From = Person{
		Email: c.conf.SenderEmail,
		Name:  c.conf.SenderName,
	}

	mailRequest.To = []Person{
		{
			Email: to,
		},
	}

	b, err := json.Marshal(&mailRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, sendGridBaseURL+sendGridSendEndpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Api-Token", c.conf.MailtrapAPIKey)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil

}

func (c *MailClient) submitSendgrid(subject, to string, content, contentType string) error {
	from := mail.NewEmail(c.conf.SenderName, c.conf.SenderEmail)
	recipient := mail.NewEmail("", to)

	p := mail.NewPersonalization()
	p.AddTos(recipient)

	mailContent := mail.NewContent(contentType, content)

	message := new(mail.SGMailV3)
	message.SetFrom(from)
	message.AddPersonalizations(p)
	message.Subject = subject
	message.AddContent(mailContent)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	res, err := c.sendgridClient.SendWithContext(ctx, message)
	if err != nil {
		return err
	}

	if isServerError(res.StatusCode) || isClientError(res.StatusCode) {
		return fmt.Errorf("could not send email: %+v", res)
	}

	return nil
}

func isClientError(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

func isServerError(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

func NewMailClient(conf *Config) *MailClient {
	httpClient := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Second * 10,
	}

	sendgridClient := sendgrid.NewSendClient(conf.SendGridAPIKey)

	return &MailClient{
		sendgridClient: sendgridClient,
		conf:           conf,
		httpClient:     httpClient,
	}
}
