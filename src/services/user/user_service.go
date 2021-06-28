package user

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
)

type UserService interface {
	Create(user *user.User) (*user.User, rest_error.RestErr)
	Delete(id uint) rest_error.RestErr
}

type userService struct {
	userRepository user2.UserRepository
}

func NewUserService(userRepository user2.UserRepository) UserService {
	return &userService{
		userRepository,
	}
}

func (s *userService) Create(user *user.User) (*user.User, rest_error.RestErr) {
	return s.userRepository.Create(user)
}

func (s *userService) Delete(id uint) rest_error.RestErr {
	return s.userRepository.Delete(id)
}
