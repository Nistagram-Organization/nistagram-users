package registered_user

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/registered_user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/mock"
)

type RegisteredUserRepositoryMock struct {
	mock.Mock
}

func (r *RegisteredUserRepositoryMock) Create(userEntity *registered_user.RegisteredUser) (*registered_user.RegisteredUser, rest_error.RestErr) {
	args := r.Called(userEntity)
	if args.Get(1) == nil {
		return args.Get(0).(*registered_user.RegisteredUser), nil
	}
	return nil, args.Get(1).(rest_error.RestErr)
}

func (r *RegisteredUserRepositoryMock) Delete(u uint) rest_error.RestErr {
	panic("implement me")
}

