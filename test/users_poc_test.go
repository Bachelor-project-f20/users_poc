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
	pb "github.com/grammeaway/users_poc/models/proto/gen"
	"github.com/grammeaway/users_poc/pkg/creating"
	"github.com/grammeaway/users_poc/pkg/deleting"
	eventHandler "github.com/grammeaway/users_poc/pkg/event"
	"github.com/grammeaway/users_poc/pkg/updating"
	"github.com/nats-io/go-nats"
)

var eventEmitter etg.EventEmitter
var eventListener etg.EventListener
var eventChan <-chan etg.Event
var creatingService creating.Service
var updatingService updating.Service
var deletingService deleting.Service

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

	eventEmitter, err = stan.NewNatsEventEmitter(encodedConn, config.Exchange, config.QueueType)

	if err != nil {
		fmt.Printf("Creation of event emitter failed, error: %v \n", err)
		t.Error(err)
	}

	outbox, obErr := outbox.NewOutbox(config.DatabaseType, config.DatabaseConnection, eventEmitter, pb.User{})

	if obErr != nil {
		fmt.Printf("Error creating Outbox: %v \n", err)
		t.Error(err)
	}

	eventListener, err = stan.NewNatsEventListener(encodedConn, config.Exchange, config.QueueType)

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

	creatingService = creating.NewService(outbox)
	updatingService = updating.NewService(outbox)
	deletingService = deleting.NewService(outbox)
}

func test(t *testing.T) {
	fmt.Println("IN TEST SETUP!!!!")
	testingChan := make(chan eventHandler.TestObject)
	defer close(testingChan)
	go func() {
		eventHandler.TestingStartEventHandler(
			testingChan,
			eventChan,
			creatingService,
			updatingService,
			deletingService,
		)
		fmt.Println("IN TEST SETUP!!!!")
	}()
	fmt.Println("HERE")
	testResult := <-testingChan
	if !testResult.Ok {
		fmt.Println("ERROR")
		t.Error(testResult.Err)
	}
	fmt.Println("OK")

	testingChan <- eventHandler.TestObject{}
}

func TestCreateRequestHandling(t *testing.T) {
	fmt.Println("TestCreateRequestHandling")
	user := &pb.User{
		ID:       "test",
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
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalUser,
	}

	eventEmitter.Emit(creationRequest)
	test(t)
}

func TestUpdateRequestHandling(t *testing.T) {
	fmt.Println("TestUpdateRequestHandling")
	user := &pb.User{
		ID:       "test",
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
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalUser,
	}

	eventEmitter.Emit(updateRequest)
	test(t)
}

func TestDeleteRequestHandling(t *testing.T) {
	fmt.Println("TestDeleteRequestHandling")
	user := &pb.User{
		ID:       "test",
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
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalUser,
	}

	eventEmitter.Emit(deletionRequest)
	test(t)
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
