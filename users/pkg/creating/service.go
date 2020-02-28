package creating

import (
	"fmt"

	pb "github.com/grammeaway/users_poc/users/models/proto/gen"
)

type Outbox interface {
	Insert(interface{}) error //ID of sorts for return val as well?
}

type Service struct {
	ob Outbox
}

func (srv *Service) CreateUser(user *pb.User) error {
	err := srv.ob.Insert(user)

	if err != nil {
		fmt.Println("Error during creation of user. Err: ", err)
	}

	//publish event here?
	return err
}
