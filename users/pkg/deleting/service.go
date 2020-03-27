package deleting

import (
	"fmt"
	"time"

	ob "github.com/dueruen/go-outbox"
	pb "github.com/grammeaway/users_poc/users/models/proto/gen"
)

type Service interface {
	DeleteUser(requestEvent ob.Event) error
}

type service struct {
	ob ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) DeleteUser(requestEvent ob.Event) error {

	deletionEvent := ob.Event{
		ID:        "test",
		Publisher: "test",
		EventName: "user_deleted",
		Timestamp: time.Now().UnixNano(),
		Payload:   []byte("test"),
	}

	userID := string(deletionEvent.Payload)

	userToDelete := &pb.User{
		ID: userID,
	}

	err := srv.ob.Delete(userToDelete, deletionEvent)

	if err != nil {
		fmt.Println("Error during deletion of user. Err: ", err)
	}

	return err
}
