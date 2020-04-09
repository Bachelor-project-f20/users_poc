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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/golang/protobuf/proto"
)

// To run the tests, a local Docker container, based on the following Docker image
// must be running, on port 9911 (unless the test code is changed):
// https://hub.docker.com/r/s12v/sns/?fbclid=IwAR23X1mEVHH5Q64awf-ZtyzC_r712-yjfmqEQGRvDCT8LYfMkdyP4goTxdE

// Alternatively, one can attach to an actual SNS instance, by using the SharedConfigState session initialization

var eventEmitter etg.EventEmitter
var eventListener etg.EventListener
var eventChan <-chan models.Event
var creatingService creating.Service
var updatingService updating.Service
var deletingService deleting.Service
var svc *sns.SNS

func TestServiceSetup(t *testing.T) {
	//AnonymousCredentials for the mock SNS instance
	//SSL disabled, because it's easier when testing
	//localhost:991 is where the fake SNS container should be running
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Credentials: credentials.AnonymousCredentials, Endpoint: aws.String("http://localhost:9911"), Region: aws.String("us-east-1"), DisableSSL: aws.Bool(true)},
	}))

	svc = sns.New(sess)

	incomingEvents := []string{
		models.UserEvents_CREATE_USER.String(),
		models.UserEvents_DELETE_USER.String(),
		models.UserEvents_UPDATE_USER.String()}

	outgoingEvents := []string{
		models.UserEvents_USER_CREATED.String(),
		models.UserEvents_USER_UPDATED.String(),
		models.UserEvents_USER_DELETED.String()}

	incomingAndOutgoingEvents := append(incomingEvents, outgoingEvents...)

	configRes, err := config.ConfigService(
		"configFile",
		config.ConfigValues{
			UseEmitter:        true,
			UseListener:       true,
			MessageBrokerType: etg.SNS,
			SNSClient:         svc,
			Events:            incomingAndOutgoingEvents,
			UseOutbox:         true,
			OutboxModels:      []interface{}{models.User{}},
		},
	)
	if err != nil {
		log.Fatalln("configuration failed, error: ", err)
		panic("configuration failed")
	}
	eventEmitter = configRes.EventEmitter
	eventListener = configRes.EventListener

	eventChan, _, err = eventListener.Listen(incomingAndOutgoingEvents...)

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

	eventEmitter.Emit(creationRequest)
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
