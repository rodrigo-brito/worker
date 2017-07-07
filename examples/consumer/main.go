package main

import (
	"fmt"

	"github.com/guilhermehubner/worker"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %s", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel: %s", err)
		return
	}
	defer ch.Close()

	ctx := context.Background()

	wp := worker.NewWorkerPool(ctx, 5, ch)

	wp.RegisterJob(worker.JobType{
		Name: "fila1",
		Handle: func(_ context.Context, msg []byte) error {
			fmt.Printf("Job: fila 1, msg: %s\n", msg)
			return nil
		},
		Priority: 10,
	})

	wp.RegisterJob(worker.JobType{
		Name: "fila2",
		Handle: func(_ context.Context, msg []byte) error {
			fmt.Printf("Job: fila 2, msg: %s\n", msg)
			return nil
		},
		Priority: 15,
	})

	wp.Start()
}
