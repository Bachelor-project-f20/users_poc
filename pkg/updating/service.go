package updating

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
	UpdateUser(requestEvent etg.Event) error
}

type service struct {
	ob ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) UpdateUser(requestEvent etg.Event) error {

	user := &pb.User{}
	err := proto.Unmarshal(requestEvent.Payload, user)

	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	// user := &pb.User{
	// 	ID:       payload.User.ID,
	// 	OfficeID: payload.User.OfficeID,
	// 	Name:     payload.User.Name,
	// }

	userUpdatedEvent := &pb.UserUpdated{
		User: user,
	}

	marshalEvent, err := proto.Marshal(userUpdatedEvent)

	if err != nil {
		log.Fatalf("Error with proto: %v \n", err)
		return err
	}

	updateEvent := etg.Event{
		ID:        "test",
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
