package notifications

import (
	"github.com/oyamo/forumz-notification-server/internal/domain/notification"
)

type Worker struct {
	JobQueue   chan *notification.Notification
	WorkerPool chan chan *notification.Notification
	updateFunc func(notification *notification.Notification) error
	sendFunc   func(notification *notification.Notification) error
}

func (w *Worker) Execute() {

	for {
		w.WorkerPool <- w.JobQueue

		select {
		case req := <-w.JobQueue:
			err := w.sendFunc(req)
			req.Sent = err == nil
			req.AttemptCount++
			err = w.updateFunc(req)
			if err != nil {
				// TODO
			}
		}
	}
}
