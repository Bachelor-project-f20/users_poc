package pkg

import (
	"log"

	"github.com/Bachelor-project-f20/shared/config"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/Bachelor-project-f20/users_poc/pkg/creating"
	"github.com/Bachelor-project-f20/users_poc/pkg/deleting"
	handler "github.com/Bachelor-project-f20/users_poc/pkg/event"
	"github.com/Bachelor-project-f20/users_poc/pkg/updating"
)

var configFile string = "configPath"

func Run() {
	configRes, err := config.ConfigService(
		configFile,
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

	incomingEvents := []string{
		models.UserEvents_CREATE_USER.String(),
		models.UserEvents_DELETE_USER.String(),
		models.UserEvents_UPDATE_USER.String()}

	eventChan, _, err := configRes.EventListener.Listen(incomingEvents...)

	if err != nil {
		log.Fatalf("Creation of subscriptions failed, error: %v \n", err)
	}

	creatingService := creating.NewService(configRes.Outbox)
	updatingService := updating.NewService(configRes.Outbox)
	deletingService := deleting.NewService(configRes.Outbox)

	handler.StartEventHandler(
		eventChan,
		creatingService,
		updatingService,
		deletingService)
}
