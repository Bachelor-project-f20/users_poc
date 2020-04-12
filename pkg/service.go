package pkg

import (
	"fmt"
	"log"
	"net/http"

	etg "github.com/Bachelor-project-f20/eventToGo"
	"github.com/Bachelor-project-f20/shared/config"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/Bachelor-project-f20/users_poc/pkg/creating"
	"github.com/Bachelor-project-f20/users_poc/pkg/deleting"
	handler "github.com/Bachelor-project-f20/users_poc/pkg/event"
	"github.com/Bachelor-project-f20/users_poc/pkg/updating"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var configFile string = "configPath"

func Run() {

	// AnonymousCredentials for the mock SNS instance
	// SSL disabled, because it's easier when testing
	// localhost:991 is where the fake SNS container should be running
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Credentials: credentials.AnonymousCredentials, Endpoint: aws.String("http://localhost:9911"), Region: aws.String("us-east-1"), DisableSSL: aws.Bool(true)},
	}))

	svc := sns.New(sess)

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

	eventChan, _, err := configRes.EventListener.Listen(incomingEvents...)

	if err != nil {
		log.Fatalf("Creation of subscriptions failed, error: %v \n", err)
	}

	creatingService := creating.NewService(configRes.Outbox)
	updatingService := updating.NewService(configRes.Outbox)
	deletingService := deleting.NewService(configRes.Outbox)

	go func() {
		fmt.Println("Serving metrics API")

		h := http.NewServeMux()
		h.Handle("/metrics", promhttp.Handler())

		http.ListenAndServe(":9191", h)
	}()

	handler.StartEventHandler(
		eventChan,
		creatingService,
		updatingService,
		deletingService)
}
