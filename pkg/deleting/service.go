package deleting

import (
	"fmt"
	"time"

	etg "github.com/Bachelor-project-f20/eventToGo"
	ob "github.com/Bachelor-project-f20/go-outbox"
	"github.com/golang/protobuf/proto"
	pb "github.com/grammeaway/users_poc/models/proto/gen"
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

	payload := &pb.DeleteUser{}

	err := proto.Unmarshal(requesEvent.GetPayload())

	if err != nil {
		return err
	}

	user := &pb.User{
		ID:       "test",
		OfficeID: payload.User.OfficeID,
		Name:     payload.User.Name,
	}

	userDeletedEvent := &pb.UserDeleted{
		User: user,
	}
	marshalEvent, err := proto.Marshal(userDeletedEvent)

	if err != nil {
		return err
	}

	deletionEvent := etg.Event{
		ID:        "test",
		Publisher: "users",
		EventName: "user_deleted",
		Timestamp: time.Now().UnixNano(),
		Payload:   marshalEvent,
	}

	err := srv.ob.Delete(userToDelete, deletionEvent)

	if err != nil {
		fmt.Println("Error during deletion of user. Err: ", err)
	}

	return err
}
