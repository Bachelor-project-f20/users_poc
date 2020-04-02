package creating

import (
	"log"
	"time"

	ob "github.com/Bachelor-project-f20/go-outbox"
	models "github.com/Bachelor-project-f20/shared/models"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/proto"
)

type Service interface {
	CreateUser(requestEvent models.Event) error
}

type service struct {
	outbox ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) CreateUser(requestEvent models.Event) error {

	event := &models.CreateUser{}
	err := proto.Unmarshal(requestEvent.Payload, event)
	user := event.User

	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	userCreatedEvent := &models.UserCreated{
		User: user,
	}
	marshalEvent, err := proto.Marshal(userCreatedEvent)
	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	id, _ := uuid.NewV4()
	idAsString := id.String()

	creationEvent := models.Event{
		ID:        idAsString,
		Publisher: models.UserService_USERS.String(),
		EventName: models.UserEvents_USER_CREATED.String(),
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalEvent,
	}

	err = srv.outbox.Insert(user, creationEvent)

	if err != nil {
		log.Fatal("Error during creation of user. Err: ", err)
	}
	return err
}
