package updating

import (
	"fmt"

	pb "github.com/grammeaway/users_poc/users/models/proto/gen"
)

type Outbox interface {
	Update(interface{}) error //ID of sorts for return val as well?
}

type Service struct {
	ob Outbox
}

func (srv *Service) UpdateUser(user *pb.User) error {
	err := srv.ob.Update(user)

	if err != nil {
		fmt.Println("Error during update of user. Err: ", err)
	}

	//publish event here?
	return err
}
