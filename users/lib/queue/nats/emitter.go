package nats

import (
	"log"

	ob "github.com/dueruen/go-outbox"
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

func (n *natsEventEmitter) Emit(e ob.Event) error {
	err := n.connection.Publish(e.EventName, e)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
