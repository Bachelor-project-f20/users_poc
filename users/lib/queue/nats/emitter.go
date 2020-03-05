package nats

import (
	"log"

	"github.com/grammeaway/users_poc/users/lib/queue"
	"github.com/nats-io/go-nats"
)

type natsEventEmitter struct {
	connection *nats.EncodedConn
	exchange   string
	queue      string
}

func NewNatsEventEmitter(connection *nats.EncodedConn, exchange, queue string) (queue.EventEmitter, error) {
	emitter := natsEventEmitter{
		connection,
		exchange,
		queue,
	}

	return &emitter, nil
}

func (n *natsEventEmitter) Emit(e queue.Event) error {
	err := n.connection.Publish(e.GetID(), e)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
