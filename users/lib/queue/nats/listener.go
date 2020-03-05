package nats

import (
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
	return &listener, nil
}

func (n *natsEventListener) Listen(events ...string) (<-chan queue.Event, <-chan error, error) {
	eventChan := make(chan queue.Event)
	errChan := make(chan error)

	for count, _ := range events {
		_, err := n.connection.QueueSubscribe(events[count], n.queue, func(e *Event) {
			eventChan <- e
		})
		if err != nil {
			errChan <- err
		}
	}
	return eventChan, errChan, nil
}
