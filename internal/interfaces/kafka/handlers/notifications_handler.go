package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/oyamo/forumz-notification-server/internal/domain/notification"
	"go.uber.org/zap"
	"path"
	"text/template"
)

type NotificationHandler struct {
	logger         *zap.SugaredLogger
	notificationUC *notification.Usecase
}

const (
	templDir       = "./web/templates"
	connectionFile = "connection.template.html"
	postsFile      = "posts.template.html"
)

func (h *NotificationHandler) parseTemplate(body map[string]string, file string) (string, error) {
	file = path.Join(templDir, file)
	tFile, err := template.ParseFiles(file)
	if err != nil {
		return "", err
	}

	writer := bytes.NewBufferString("")
	err = tFile.Execute(writer, body)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}

func (h *NotificationHandler) HandleNotification(b []byte) {
	var dto notification.NotificationDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		h.logger.Errorw("error unmarshalling notification", "error", err)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		h.logger.Errorw("error generating id", "error", err)
		return
	}

	var content string
	var priority int

	switch dto.Type {
	case notification.NotificationTypeConnection:
		content, err = h.parseTemplate(dto.AdditionalInfo, connectionFile)
		priority = 2
	case notification.NotificationTypePost:
		content, err = h.parseTemplate(dto.AdditionalInfo, postsFile)
		priority = 3
	default:
		err = fmt.Errorf("unknown notification type: %s", dto.Type)
	}

	if err != nil {
		h.logger.Errorw("error parsing template", "error", err)
		return
	}

	not := notification.Notification{
		ID:           id,
		EmailAddress: dto.Recipient,
		Content:      content,
		Type:         dto.Type,
		Priority:     priority,
	}

	ctx := context.Background()
	err = h.notificationUC.Save(ctx, &not)
	if err != nil {
		h.logger.Errorw("error saving notification", "error", err)
	}

}

func NewNotificationHandler(logger *zap.SugaredLogger,
	notificationUC *notification.Usecase) *NotificationHandler {
	return &NotificationHandler{
		logger:         logger,
		notificationUC: notificationUC,
	}
}
