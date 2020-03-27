package creating

import (
	"fmt"
	"time"

	ob "github.com/Bachelor-project-f20/go-outbox"
	pb "github.com/grammeaway/users_poc/users/models/proto/gen"
)

type Service interface {
	CreateUser(requestEvent ob.Event) error
}

type service struct {
	ob ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) CreateUser(requestEvent ob.Event) error {

	creationEvent := ob.Event{
		ID:        "test",
		Publisher: "test",
		EventName: "user_created",
		Timestamp: time.Now().UnixNano(),
		Payload:   []byte("test"),
	}

	user := &pb.User{
		ID:       "test",
		OfficeID: "test",
		Name:     "test",
	}

	err := srv.ob.Insert(user, creationEvent)

	if err != nil {
		fmt.Println("Error during creation of user. Err: ", err)
	}
	return err
}
