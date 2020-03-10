package updating

import (
	"fmt"
	"time"

	ob "github.com/dueruen/go-outbox"
	pb "github.com/grammeaway/users_poc/users/models/proto/gen"
)

type Service interface {
	UpdateUser(user *pb.User) error
}

type service struct {
	ob ob.Outbox
}

func NewService(outbox ob.Outbox) Service {
	return &service{outbox}
}

func (srv *service) UpdateUser(user *pb.User) error {

	event := ob.Event{
		ID:        "test",
		Publisher: "test",
		EventName: "user_updated",
		Timestamp: time.Now().UnixNano(),
		Payload:   []byte("test"),
	}

	err := srv.ob.Update(user, event)

	if err != nil {
		fmt.Println("Error during update of user. Err: ", err)
	}

	return err
}
