package nats

import (
	"fmt"

	"github.com/grammeaway/users_poc/users/lib/queue"
	"github.com/nats-io/go-nats"
)

type natsEventListener struct {
	connection *nats.EncodedConn
	exchange   string
	queue      string
}

type Event struct {
	ID        string
	Publisher string
	Timestamp int64
}

func (e *Event) GetID() string {
	return e.ID
}

func (e *Event) GetPublisher() string {
	return e.Publisher
}

func (e *Event) GetTimestamp() int64 {
	return e.Timestamp
}

func NewNatsEventListener(connection *nats.EncodedConn, exchange, queue string) (queue.EventListener, error) {
	listener := natsEventListener{
		connection,
		exchange,
		queue,
	}

	// err := listener.setup()
	// if err != nil {
	// 	return nil, err
	// }

	return &listener, nil
}

// func (n *natsEventListener) setup() error {
// 	eventChan := make(chan queue.Event)
// 	errChan := make(chan error)

// 	if err != nil {
// 		fmt.Println("Error creating message queue subscription")
// 		return err
// 	}

// 	return nil
// }

func (n *natsEventListener) Listen(events ...string) (<-chan queue.Event, <-chan error, error) {
	fmt.Println("Listen invoked")
	eventChan := make(chan queue.Event)
	errChan := make(chan error)

	for count, _ := range events {
		fmt.Println("In the loop")
		//fmt.Println(events[count])
		//sub, err := n.connection.QueueSubscribe(events[count], n.queue, func(e *queue.Event) {
		sub, err := n.connection.QueueSubscribe(events[count], n.queue, func(e *Event) {
			fmt.Printf("Event received: %v \n", e)
			eventChan <- e
		})
		fmt.Printf("Subscription established: %v \n", sub.Subject)

		if err != nil {
			errChan <- err
		}
	}

	return eventChan, errChan, nil
}
