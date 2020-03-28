package test

import (
	"fmt"
	"testing"
	"time"

	etg "github.com/Bachelor-project-f20/eventToGo"
	stan "github.com/Bachelor-project-f20/eventToGo/nats"
	"github.com/Bachelor-project-f20/go-outbox"
	"github.com/golang/protobuf/proto"
	"github.com/grammeaway/users_poc/lib/configure"
	"github.com/grammeaway/users_poc/pkg/creating"
	"github.com/grammeaway/users_poc/pkg/deleting"
	"github.com/grammeaway/users_poc/pkg/event/handler"
	"github.com/grammeaway/users_poc/pkg/updating"
	"github.com/nats-io/go-nats"
)

var eventEmitter stan.eventEmitter
var eventListener stan.eventListener
var eventHandler handler.Service
var eventChan <-chan etg.Event

func TestServiceSetup(t *testing.T) {
	configFile := "testConfig"
	config, err := configure.ExtractConfiguration(configFile)

	if err != nil {
		fmt.Printf("Error extracting config file: %v \n", err)
		fmt.Println("Using default configuration")

	}

	encodedConn, err := setupNatsConn()

	if err != nil {
		fmt.Printf("Error connecting to Nats: %v \n", err)
		t.Error(err)
	}

	exchange := "test"
	queueType := "queue"

	eventEmitter, err = stan.NewNatsEventEmitter(encodedConn, exchange, queueType)

	if err != nil {
		fmt.Printf("Creation of event emitter failed, error: %v \n", err)
		t.Error(err)
	}

	_, obErr := outbox.NewOutbox(config.DatabaseType, config.DatabaseConnection, eventEmitter)

	if obErr != nil {
		fmt.Printf("Error creating Outbox: %v \n", err)
		t.Error(err)
	}

	eventListener, err = stan.NewNatsEventListener(encodedConn, exchange, queueType)

	if err != nil {
		fmt.Printf("Creation of Listener  failed, error: %v", err)
		t.Error(err)
	}

	incomingEvents := []string{"creation_request", "updating_request", "deletion_request"} //I'm guessing this should probably go in the proto files?

	eventChan, _, err = eventListener.Listen(incomingEvents...)

	if err != nil {
		fmt.Printf("Creation of subscriptions failed, error: %v \n", err)
		t.Error(err)
	}

	creatingService := creating.NewService(outbox)
	updatingService := updating.NewService(outbox)
	deletingService := deleting.NewService(outbox)

	eventHandler = handler.NewEventHandler(creatingService, updatingService, deletingService)
}

func TestCreateRequestHandling(t *testing.T) {

	user := &pb.User{
		ID:       "test", //TODO, make ID
		OfficeID: "test",
		Name:     "creation_test_user",
	}

	marshalUser, err := proto.Marshal(user)

	if err != nil {
		fmt.Printf("Error marshalling new user, error: %v \n", err)
		t.Error(err)
	}

	creationRequest := etg.Event{
		ID:        "test",
		Publisher: "users_test",
		EventName: "creation_request",
		TimeStamp: time.Now().UnixNano(),
		Payload:   marshalUser,
	}

	eventEmitter.Emit(creationRequest)
	//How to check if this actually works?

}

func TestUpdateRequestHandling(t *testing.T) {
	user := &pb.User{
		ID:       "test", //TODO, make ID
		OfficeID: "new_office_id",
		Name:     "creation_test_user",
	}

	marshalUser, err := proto.Marshal(user)

	if err != nil {
		fmt.Printf("Error marshalling new user, error: %v \n", err)
		t.Error(err)
	}

	updateRequest := etg.Event{
		ID:        "test",
		Publisher: "users_test",
		EventName: "updating_request",
		TimeStamp: time.Now().UnixNano(),
		Payload:   marshalUser,
	}

	eventEmitter.Emit(updateRequest)
	//How to check if this actually works?
}

func TestDeleteRequestHandling(t *testing.T) {
	user := &pb.User{
		ID:       "test", //TODO, make ID
		OfficeID: "new_office_id",
		Name:     "creation_test_user",
	}

	marshalUser, err := proto.Marshal(user)

	if err != nil {
		fmt.Printf("Error marshalling new user, error: %v \n", err)
		t.Error(err)
	}

	deletionRequest := etg.Event{
		ID:        "test",
		Publisher: "users_test",
		EventName: "deletion_request",
		TimeStamp: time.Now().UnixNano(),
		Payload:   marshalUser,
	}

	eventEmitter.Emit(deletionRequest)

	//How to check if this actually works?
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
