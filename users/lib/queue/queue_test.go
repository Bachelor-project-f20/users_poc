package queue_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/grammeaway/users_poc/users/lib/queue"
	lnats "github.com/grammeaway/users_poc/users/lib/queue/nats"
	"github.com/nats-io/go-nats"
)

var natsConn *nats.Conn
var encodedConn *nats.EncodedConn
var eventListener queue.EventListener
var eventEmitter queue.EventEmitter

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

func TestEmitAndListen(t *testing.T) {
	natsConn, err := nats.Connect("localhost:4222")

	if err != nil {
		fmt.Println("Connection to Nats failed")
		t.Error(err)
	}

	encodedConn, err := nats.NewEncodedConn(natsConn, "json")

	if err != nil {
		fmt.Println("Creation of encoded connection failed")
		t.Error(err)
	}

	exchange := "test"
	eventQueue := "queue"

	eventEmitter, err := lnats.NewNatsEventEmitter(encodedConn, exchange, eventQueue)

	if err != nil {
		fmt.Println("Creation of Emitter  failed")
		t.Error(err)
	}

	eventListener, err := lnats.NewNatsEventListener(encodedConn, exchange, eventQueue)

	if err != nil {
		fmt.Println("Creation of Listener  failed")
		t.Error(err)
	}

	event := &Event{
		"test",
		"test",
		time.Now().UnixNano(),
	}

	eventChan, _, err := eventListener.Listen(event.GetID())
	if err != nil {
		fmt.Println("Listen function  failed")
		t.Error(err)
	}

	eventEmitter.Emit(event)
	fmt.Println("Event emited")
	recEvent := <-eventChan

	fmt.Printf("Event received in test: %v", recEvent)
}
