package updating

import (
	"log"
	"time"

	ob "github.com/Bachelor-project-f20/go-outbox"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/proto"
)

type Service interface {
	UpdateUser(requestEvent models.Event) error
}

type service struct {
	ob ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) UpdateUser(requestEvent models.Event) error {

	event := &models.CreateUser{}
	err := proto.Unmarshal(requestEvent.Payload, event)
	user := event.User

	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	userUpdatedEvent := &models.UserUpdated{
		User: user,
	}

	marshalEvent, err := proto.Marshal(userUpdatedEvent)

	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	id, _ := uuid.NewV4()
	idAsString := id.String()

	updateEvent := models.Event{
		ID:        idAsString,
		Publisher: models.UserService_USERS.String(),
		EventName: models.UserEvents_USER_UPDATED.String(),
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalEvent,
		ApiTag:    requestEvent.ApiTag,
	}

	err = srv.ob.Update(user, updateEvent)

	if err != nil {
		log.Fatal("Error during update of user. Err: ", err)
	}

	return err
}
