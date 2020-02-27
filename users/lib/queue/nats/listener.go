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

func NewNatsEventListener(connection *nats.EncodedConn, exchange, queue string) (queue.EventListener, error) {
	listener := natsEventListener { 
		connection,
		exchange,
		queue
	}

	err := listener.setup()
	if err != nil { 
		return nil, err
	}

	return &listener, nil
}

func (n *natsEventListener) setup() error { 
	channel, err := n.connecti
}
