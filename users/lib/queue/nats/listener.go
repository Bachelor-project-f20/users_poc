package nats

import (
	ob "github.com/dueruen/go-outbox"
	"github.com/grammeaway/users_poc/users/lib/queue"
	"github.com/nats-io/go-nats"
)

type natsEventListener struct {
	connection *nats.EncodedConn
	exchange   string
	queue      string
}

func NewNatsEventListener(connection *nats.EncodedConn, exchange, queue string) (queue.EventListener, error) {
	listener := natsEventListener{
		connection,
		exchange,
		queue,
	}
	return &listener, nil
}

func (n *natsEventListener) Listen(events ...string) (<-chan ob.Event, <-chan error, error) {
	eventChan := make(chan ob.Event)
	errChan := make(chan error)

	for count, _ := range events {
		_, err := n.connection.QueueSubscribe(events[count], n.queue, func(e *ob.Event) {
			eventChan <- *e
		})
		if err != nil {
			errChan <- err
		}
	}
	return eventChan, errChan, nil
}
