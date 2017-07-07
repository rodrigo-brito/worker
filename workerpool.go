package worker

import (
	"sort"

	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

type workerPool struct {
	ctx      context.Context
	channel  *amqp.Channel
	workers  []*worker
	jobTypes jobTypes
	cancel   chan bool
}

// Start starts the workers and associated processes.
func (wp *workerPool) Start() {
	sort.Sort(wp.jobTypes)

	getJob := func() (*amqp.Delivery, jobHandle, uint8) {
		for _, jobType := range wp.jobTypes {
			msg, ok, _ := wp.channel.Get(jobType.Name, false)
			if !ok {
				continue
			}

			return &msg, jobType.Handle, jobType.Retry
		}

		return nil, nil, 0
	}

	for _, w := range wp.workers {
		w = newWorker(wp.ctx, getJob, wp.cancel)
		w.start()
	}

	<-wp.cancel
}

// RegisterJob adds a job with handler for 'name' queue and allows you to specify options such as a job's priority and it's retry count.
func (wp *workerPool) RegisterJob(job JobType) {
	wp.jobTypes = append(wp.jobTypes, job)
}

// NewWorkerPool creates a new worker pool. ctx will be used for middleware and handlers. concurrency specifies how many workers to spin up - each worker can process jobs concurrently.
func NewWorkerPool(ctx context.Context, concurrency uint, channel *amqp.Channel) *workerPool {
	if channel == nil {
		panic("worker equeuer: needs a non-nil *amqp.Channel")
	}

	wp := &workerPool{
		ctx:     ctx,
		channel: channel,
		workers: make([]*worker, concurrency),
	}

	return wp
}
