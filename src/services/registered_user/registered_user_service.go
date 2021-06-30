package registered_user

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/registered_user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	regUserRepo "github.com/Nistagram-Organization/nistagram-users/src/repositories/registered_user"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
)

type RegisteredUserService interface {
	Create(*registered_user.RegisteredUser) (*registered_user.RegisteredUser, rest_error.RestErr)
	Delete(uint) rest_error.RestErr
}

type registeredUserService struct {
	registeredUserRepository regUserRepo.RegisteredUserRepository
	userRepository           user.UserRepository
}

func NewRegisteredUserService(registeredUserRepository regUserRepo.RegisteredUserRepository, userRepository user.UserRepository) RegisteredUserService {
	return &registeredUserService{
		registeredUserRepository,
		userRepository,
	}
}

func (s *registeredUserService) Create(user *registered_user.RegisteredUser) (*registered_user.RegisteredUser, rest_error.RestErr) {
	return s.registeredUserRepository.Create(user)
}

func (s *registeredUserService) Delete(id uint) rest_error.RestErr {
	if delErr := s.userRepository.Delete(id); delErr != nil {
		return delErr
	}

	return s.registeredUserRepository.Delete(id)
}
