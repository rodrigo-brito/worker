package worker

import (
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type getJobHandle func() (*amqp.Delivery, jobHandle, uint8)

type worker struct {
	ctx    context.Context
	getJob getJobHandle
	cancel chan bool
}

func newWorker(ctx context.Context, getJob getJobHandle,
	cancel chan bool) *worker {
	return &worker{
		ctx:    ctx,
		getJob: getJob,
		cancel: cancel,
	}
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case <-w.cancel:
				return
			default:
				w.executeJob()
			}
		}
	}()
}

func (w *worker) executeJob() {
	message, handle, retries := w.getJob()
	if message == nil || handle == nil {
		return
	}
	if retries == 0 {
		retries = 1
	}

	for i := retries; i > 0; i-- {
		err := handle(w.ctx, message.Body)
		if err == nil {
			break
		}
	}

	message.Ack(false)
}
