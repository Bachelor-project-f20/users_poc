package creating

import (
	"fmt"
	"log"
	"time"

	etg "github.com/Bachelor-project-f20/eventToGo"
	ob "github.com/Bachelor-project-f20/go-outbox"
	"github.com/golang/protobuf/proto"
	pb "github.com/grammeaway/users_poc/models/proto/gen"
)

type Service interface {
	CreateUser(requestEvent etg.Event) error
}

type service struct {
	outbox ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) CreateUser(requestEvent etg.Event) error {

	user := &pb.User{}
	err := proto.Unmarshal(requestEvent.Payload, user)

	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	// user := &pb.User{
	// 	ID:       "test", //TODO, make ID
	// 	OfficeID: payload.User.OfficeID,
	// 	Name:     payload.User.Name,
	// }

	userCreatedEvent := &pb.UserCreated{
		User: user,
	}
	marshalEvent, err := proto.Marshal(userCreatedEvent)
	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	creationEvent := etg.Event{
		ID:        "test", //TODO, make ID
		Publisher: "users",
		EventName: "user_created",
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalEvent,
	}

	err = srv.outbox.Insert(user, creationEvent)

	if err != nil {
		fmt.Println("Error during creation of user. Err: ", err)
	}
	return err
}
