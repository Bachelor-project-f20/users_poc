package updating

import (
	"fmt"
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

	user := &models.User{}
	err := proto.Unmarshal(requestEvent.Payload, user)

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
		Publisher: "users",
		EventName: "user_updated",
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalEvent,
	}

	err = srv.ob.Update(user, updateEvent)

	if err != nil {
		fmt.Println("Error during update of user. Err: ", err)
	}

	return err
}
