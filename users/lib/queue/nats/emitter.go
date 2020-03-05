package nats

import (
	"fmt"
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
	fmt.Printf("Event emitted for subject: %s \n", e.GetID())
	if err != nil {
		log.Fatal(err)
		return err
	}
	//fmt.Println("Event emitted")
	return nil
}
