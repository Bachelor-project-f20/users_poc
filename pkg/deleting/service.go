package deleting

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
	DeleteUser(requestEvent models.Event) error
}

type service struct {
	ob ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) DeleteUser(requestEvent models.Event) error {

	user := &models.User{}

	err := proto.Unmarshal(requestEvent.Payload, user)

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
		Publisher: "users",
		EventName: "user_deleted",
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalEvent,
	}

	err = srv.ob.Delete(user, deletionEvent)

	if err != nil {
		fmt.Println("Error during deletion of user. Err: ", err)
	}

	return err
}
