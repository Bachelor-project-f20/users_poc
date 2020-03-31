package deleting

import (
	"fmt"
	"log"
	"time"

	etg "github.com/Bachelor-project-f20/eventToGo"
	ob "github.com/Bachelor-project-f20/go-outbox"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/proto"
	pb "github.com/Bachelor-project-f20/users_poc/models/proto/gen"
)

type Service interface {
	DeleteUser(requestEvent etg.Event) error
}

type service struct {
	ob ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) DeleteUser(requestEvent etg.Event) error {

	user := &pb.User{}

	err := proto.Unmarshal(requestEvent.Payload, user)

	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	userDeletedEvent := &pb.UserDeleted{
		User: user,
	}
	marshalEvent, err := proto.Marshal(userDeletedEvent)

	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	id, _ := uuid.NewV4()
	idAsString := id.String()

	deletionEvent := etg.Event{
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
