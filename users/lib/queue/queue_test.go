package queue_test

import (
	"fmt"
	"testing"
	"time"

	lnats "github.com/grammeaway/users_poc/users/lib/queue/nats"
	"github.com/nats-io/go-nats"
)

var natsConn *nats.Conn
var encodedConn *nats.EncodedConn
var exchange string
var queueType string

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

func TestSetup(t *testing.T) {
	natsConn, err := nats.Connect("localhost:4222")

	if err != nil {
		fmt.Println("Connection to Nats failed")
		t.Error(err)
	}

	encodedConn, err = nats.NewEncodedConn(natsConn, "json")

	if err != nil {
		fmt.Println("Creation of encoded connection failed")
		t.Error(err)
	}

	exchange = "test"
	queueType = "queue"

}

func TestEmit(t *testing.T) {
	eventEmitter, err := lnats.NewNatsEventEmitter(encodedConn, exchange, queueType)

	if err != nil {
		fmt.Println("Creation of Emitter  failed")
		t.Error(err)
	}

	event := &Event{
		"test",
		"test",
		time.Now().UnixNano(),
	}

	emitErr := eventEmitter.Emit(event)

	if emitErr != nil {
		fmt.Println("Error while emitting event")
		t.Error(err)
	}
	fmt.Println("Event emited")
}

func TestListen(t *testing.T) {
	eventListener, err := lnats.NewNatsEventListener(encodedConn, exchange, queueType)

	if err != nil {
		fmt.Println("Creation of Listener  failed")
		t.Error(err)
	}

	event := &Event{
		"test",
		"test",
		time.Now().UnixNano(),
	}

	eventChan, _, listenErr := eventListener.Listen(event.GetID())
	if listenErr != nil {
		fmt.Println("Listen function  failed")
		t.Error(err)
	}

	//Necessary - when the Nats connection is not set to durable, messages in unsubscribed message queues are lost
	eventEmitter, _ := lnats.NewNatsEventEmitter(encodedConn, exchange, queueType)
	eventEmitter.Emit(event)

	recEvent := <-eventChan

	fmt.Printf("Event received in test: %v", recEvent)
}
