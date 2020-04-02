package deleting

import (
	"log"
	"time"

	ob "github.com/Bachelor-project-f20/go-outbox"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/proto"
)

type Service interface {
	DeleteUser(requestEvent models.Event) error
}

type service struct {
	ob ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) DeleteUser(requestEvent models.Event) error {

	event := &models.CreateUser{}
	err := proto.Unmarshal(requestEvent.Payload, event)
	user := event.User

	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	userDeletedEvent := &models.UserDeleted{
		User: user,
	}
	marshalEvent, err := proto.Marshal(userDeletedEvent)

	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	id, _ := uuid.NewV4()
	idAsString := id.String()

	deletionEvent := models.Event{
		ID:        idAsString,
		Publisher: models.UserService_USERS.String(),
		EventName: models.UserEvents_USER_DELETED.String(),
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalEvent,
	}

	err = srv.ob.Delete(user, deletionEvent)

	if err != nil {
		log.Fatal("Error during deletion of user. Err: ", err)
	}

	return err
}
