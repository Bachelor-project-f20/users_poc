package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	etg "github.com/Bachelor-project-f20/eventToGo"
	"github.com/Bachelor-project-f20/shared/config"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/Bachelor-project-f20/users_poc/pkg/creating"
	"github.com/Bachelor-project-f20/users_poc/pkg/deleting"
	eventHandler "github.com/Bachelor-project-f20/users_poc/pkg/event"
	"github.com/Bachelor-project-f20/users_poc/pkg/updating"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

var eventEmitter etg.EventEmitter
var eventListener etg.EventListener
var eventChan <-chan models.Event
var creatingService creating.Service
var updatingService updating.Service
var deletingService deleting.Service

func TestServiceSetup(t *testing.T) {
	configRes, err := config.ConfigService(
		"configFile",
		config.ConfigValues{
			UseEmitter:   true,
			UseListener:  true,
			UseOutbox:    true,
			OutboxModels: []interface{}{models.User{}},
		},
	)
	if err != nil {
		log.Fatalln("configuration failed, error: ", err)
		panic("configuration failed")
	}
	eventEmitter = configRes.EventEmitter
	eventListener = configRes.EventListener

	incomingEvents := []string{
		models.UserEvents_CREATE_USER.String(),
		models.UserEvents_DELETE_USER.String(),
		models.UserEvents_UPDATE_USER.String()}

	eventChan, _, err = eventListener.Listen(incomingEvents...)

	if err != nil {
		fmt.Printf("Creation of subscriptions failed, error: %v \n", err)
		t.Error(err)
	}

	creatingService = creating.NewService(configRes.Outbox)
	updatingService = updating.NewService(configRes.Outbox)
	deletingService = deleting.NewService(configRes.Outbox)
}

func test(t *testing.T) {
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
	}()
	testResult := <-testingChan
	if !testResult.Ok {
		fmt.Println("ERROR")
		t.Error(testResult.Err)
	}
	testingChan <- eventHandler.TestObject{}
}

func TestCreateRequestHandling(t *testing.T) {
	fmt.Println("TestCreateRequestHandling")
	event := models.UserCreated{
		User: &models.User{
			ID:       "test",
			OfficeID: "test",
			Name:     "creation_test_user",
		},
	}

	marshalEvent, err := proto.Marshal(&event)
	fmt.Println("HERE 1")
	if err != nil {
		fmt.Printf("Error marshalling new user, error: %v \n", err)
		t.Error(err)
	}

	creationRequest := models.Event{
		ID:        "test",
		Publisher: "users_test",
		EventName: models.UserEvents_CREATE_USER.String(),
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalEvent,
	}
	fmt.Println("HERE 2")
	eventEmitter.Emit(creationRequest)
	fmt.Println("HERE 3")
	test(t)
}

func TestUpdateRequestHandling(t *testing.T) {
	fmt.Println("TestUpdateRequestHandling")
	event := models.UserCreated{
		User: &models.User{
			ID:       "test",
			OfficeID: "new_office_id",
			Name:     "creation_test_user",
		},
	}

	marshalEvent, err := proto.Marshal(&event)

	if err != nil {
		fmt.Printf("Error marshalling new user, error: %v \n", err)
		t.Error(err)
	}

	updateRequest := models.Event{
		ID:        "test",
		Publisher: "users_test",
		EventName: models.UserEvents_UPDATE_USER.String(),
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalEvent,
	}

	eventEmitter.Emit(updateRequest)
	test(t)
}

func TestDeleteRequestHandling(t *testing.T) {
	fmt.Println("TestDeleteRequestHandling")
	event := models.UserCreated{
		User: &models.User{
			ID:       "test",
			OfficeID: "new_office_id",
			Name:     "creation_test_user",
		},
	}

	marshalEvent, err := proto.Marshal(&event)

	if err != nil {
		fmt.Printf("Error marshalling new user, error: %v \n", err)
		t.Error(err)
	}

	deletionRequest := models.Event{
		ID:        "test",
		Publisher: "users_test",
		EventName: models.UserEvents_DELETE_USER.String(),
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalEvent,
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
