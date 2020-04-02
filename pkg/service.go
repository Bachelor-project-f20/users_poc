package pkg

import (
	"fmt"
	"log"

	stan "github.com/Bachelor-project-f20/eventToGo/nats"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/Bachelor-project-f20/users_poc/lib/configure"
	"github.com/Bachelor-project-f20/users_poc/pkg/creating"
	"github.com/Bachelor-project-f20/users_poc/pkg/deleting"
	handler "github.com/Bachelor-project-f20/users_poc/pkg/event"
	"github.com/Bachelor-project-f20/users_poc/pkg/updating"
	"github.com/nats-io/nats.go"

	"github.com/Bachelor-project-f20/go-outbox"
)

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

	eventEmitter, err := stan.NewNatsEventEmitter(encodedConn, config.Exchange, config.QueueType)

	if err != nil {
		log.Fatalf("Error creating Emitter: %v \n", err)
	}

	outbox, err := outbox.NewOutbox(config.DatabaseType, config.DatabaseConnection, eventEmitter, models.User{})

	if err != nil {
		log.Fatalf("Error creating Outbox: %v \n", err)
	}

	eventListener, err := stan.NewNatsEventListener(encodedConn, config.Exchange, config.QueueType)

	if err != nil {
		log.Fatalf("Creation of Listener  failed, error: %v", err)
	}

	incomingEvents := []string{
		models.UserEvents_CREATE_USER.String(),
		models.UserEvents_DELETE_USER.String(),
		models.UserEvents_UPDATE_USER.String()}

	eventChan, _, err := eventListener.Listen(incomingEvents...)

	if err != nil {
		log.Fatalf("Creation of subscriptions failed, error: %v \n", err)
	}

	creatingService := creating.NewService(outbox)
	updatingService := updating.NewService(outbox)
	deletingService := deleting.NewService(outbox)

	handler.StartEventHandler(
		eventChan,
		creatingService,
		updatingService,
		deletingService)
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
