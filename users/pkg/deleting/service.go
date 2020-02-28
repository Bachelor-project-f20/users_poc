package deleting

import (
	"fmt"

	pb "github.com/grammeaway/users_poc/users/models/proto/gen"
)

type Outbox interface {
	Delete(interface{}) error //ID of sorts for return val as well?
}

type Service struct {
	ob Outbox
}

func (srv *Service) DeleteUser(user *pb.User) error {
	err := srv.ob.Delete(user)

	if err != nil {
		fmt.Println("Error during deletion of user. Err: ", err)
	}

	//publish event here?
	return err
}
