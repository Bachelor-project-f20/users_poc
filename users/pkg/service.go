package pkg

import (
	"fmt"
	"log"

	"github.com/grammeaway/users_poc/users/lib/queue/configure"
	stan "github.com/grammeaway/users_poc/users/lib/queue/nats"
	"github.com/grammeaway/users_poc/users/pkg/creating"
	"github.com/grammeaway/users_poc/users/pkg/deleting"
	"github.com/grammeaway/users_poc/users/pkg/event/handler"
	"github.com/grammeaway/users_poc/users/pkg/updating"
	"github.com/nats-io/go-nats"

	"github.com/dueruen/go-outbox"
)

type creatingService = creating.Service
type updatingService = updating.Service
type deletingService = deleting.Service

var configFile string = "configPath"

func Run() {

	config, err := configure.ExtractConfiguration(configFile)

	if err != nil {
		fmt.Printf("Error extracting config file: %v \n", err)
		fmt.Println("Using default configuration")
	}

	encodedConn, err := setupNatsConn()

	if err != nil {
		log.Fatalf("Error connecting to Nats: %v \n", err)
	}

	exchange := "test"
	queueType := "queue"

	eventEmitter, err := stan.NewNatsEventEmitter(encodedConn, exchange, queueType)

	if err != nil {
		log.Fatalf("Error creating Emitter: %v \n", err)
	}

	outbox, err := outbox.NewOutbox(config.DatabaseType, config.DatabaseConnection, eventEmitter)

	if err != nil {
		log.Fatalf("Error creating Outbox: %v \n", err)
	}

	eventListener, err := stan.NewNatsEventListener(encodedConn, exchange, queueType)

	if err != nil {
		log.Fatalf("Creation of Listener  failed, error: %v", err)
	}

	incomingEvents := []string{"creation_request", "updating_request", "deletion_request"} //I'm guessing this should probably go in the proto files?

	eventChan, _, err := eventListener.Listen(incomingEvents...)

	if err != nil {
		log.Fatalf("Creation of subscriptions failed, error: %v \n", err)
	}

	creatingService := creating.NewService(outbox)
	updatingService := updating.NewService(outbox)
	deletingService := deleting.NewService(outbox)

	handler := handler.NewEventHandler(creatingService, updatingService, deletingService)

	go func() {
		for {
			event := <-eventChan
			handler.HandleEvent(event)
		}
	}	

}

func setupNatsConn() (*nats.EncodedConn, error) {

	natsConn, err := nats.Connect("localhost:4222")

	if err != nil {
		fmt.Println("Connection to Nats failed")
		return nil, err
	}

	encodedConn, err := nats.NewEncodedConn(natsConn, "json")

	if err != nil {
		fmt.Println("Creation of encoded connection failed")
		return nil, err
	}

	return encodedConn, nil

}
