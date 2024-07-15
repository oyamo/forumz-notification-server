package notifications

import (
	"github.com/oyamo/forumz-notification-server/internal/domain/notification"
)

type notificationFunc func(notification *notification.Notification) error

type Dispatcher struct {
	workerPool chan chan *notification.Notification
	MaxWorkers int
	JobQueue   chan *notification.Notification
	updateFunc notificationFunc
	sendFunc   notificationFunc
}

func (d *Dispatcher) Start() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := &Worker{
			JobQueue:   make(chan *notification.Notification),
			updateFunc: d.updateFunc,
			sendFunc:   d.sendFunc,
			WorkerPool: d.workerPool,
		}
		go worker.Execute()
	}

	d.Dispatch()
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			go func(job *notification.Notification) {
				jobChannel := <-d.workerPool
				jobChannel <- job
			}(job)
		}
	}
}

func (d *Dispatcher) AddToWorkerPool(notification *notification.Notification) {
	d.JobQueue <- notification
}

func NewDispatcher(maxWorkers int, sendFunc, updateFunc notificationFunc) *Dispatcher {
	workerPool := make(chan chan *notification.Notification, maxWorkers)
	jobQueue := make(chan *notification.Notification, maxWorkers)
	return &Dispatcher{
		workerPool: workerPool,
		MaxWorkers: maxWorkers,
		JobQueue:   jobQueue,
		sendFunc:   sendFunc,
		updateFunc: updateFunc,
	}
}
