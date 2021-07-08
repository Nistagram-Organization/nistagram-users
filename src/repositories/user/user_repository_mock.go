package user

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) GetByEmail(s string) (*user.User, rest_error.RestErr) {
	args := u.Called(s)
	if args.Get(1) == nil {
		return args.Get(0).(*user.User), nil
	}
	return nil, args.Get(1).(rest_error.RestErr)
}

func (u *UserRepositoryMock) GetByUsername(s string) (*user.User, rest_error.RestErr) {
	args := u.Called(s)
	if args.Get(1) == nil {
		return args.Get(0).(*user.User), nil
	}
	return nil, args.Get(1).(rest_error.RestErr)
}

func (u *UserRepositoryMock) Update(userEntity *user.User) (*user.User, rest_error.RestErr) {
	args := u.Called(userEntity)
	if args.Get(1) == nil {
		return args.Get(0).(*user.User), nil
	}
	return nil, args.Get(1).(rest_error.RestErr)
}

func (u *UserRepositoryMock) Delete(u2 uint) rest_error.RestErr {
	panic("implement me")
}

func (u *UserRepositoryMock) DeleteFavorite(u3 uint, u2 uint) rest_error.RestErr {
	panic("implement me")
}

func (u *UserRepositoryMock) GetById(u2 uint) (*user.User, rest_error.RestErr) {
	panic("implement me")
}


