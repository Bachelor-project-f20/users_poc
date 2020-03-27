package updating

import (
	"fmt"
	"time"

	ob "github.com/Bachelor-project-f20/go-outbox"
	pb "github.com/grammeaway/users_poc/users/models/proto/gen"
)

type Service interface {
	UpdateUser(requestEvent ob.Event) error
}

type service struct {
	ob ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) UpdateUser(requestEvent ob.Event) error {

	updateEvent := ob.Event{
		ID:        "test",
		Publisher: "test",
		EventName: "user_updated",
		Timestamp: time.Now().UnixNano(),
		Payload:   []byte("test"),
	}

	userID := string(updateEvent.Payload)

	userToBeUpdated := &pb.User{
		ID: userID,
	}

	err := srv.ob.Update(userToBeUpdated, updateEvent)

	if err != nil {
		fmt.Println("Error during update of user. Err: ", err)
	}

	return err
}
