package user

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
	"net/mail"
	"net/url"
	"strings"
)

type UserService interface {
	Create(user *user.User) (*user.User, rest_error.RestErr)
	Delete(id uint) rest_error.RestErr
	Edit(user *user.User) (*user.User, rest_error.RestErr)
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

func (s *userService) Edit(user *user.User) (*user.User, rest_error.RestErr) {
	if err := ValidateForEdit(user); err != nil {
		return nil, err
	}
	existingUser, err := s.userRepository.GetByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	existingUser.Name = user.Name
	existingUser.Surname = user.Surname
	existingUser.Website = user.Website
	existingUser.Phone = user.Phone
	existingUser.Public = user.Public
	existingUser.Taggable = user.Taggable
	existingUser.Biography = user.Biography

	editedUser, err := s.userRepository.Edit(existingUser)
	if err != nil {
		return nil, err
	}

	return editedUser, nil
}

func ValidateForEdit(user *user.User) rest_error.RestErr {
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return rest_error.NewBadRequestError("Invalid email address")
	}
	if strings.TrimSpace(user.Name) == "" {
		return rest_error.NewBadRequestError("Name cannot be empty")
	}
	if strings.TrimSpace(user.Surname) == "" {
		return rest_error.NewBadRequestError("Surname cannot be empty")
	}
	if strings.TrimSpace(user.Phone) == "" {
		return rest_error.NewBadRequestError("Phone cannot be empty")
	}
	if user.Website != "" {
		if _, err := url.ParseRequestURI(user.Website); err != nil {
			return rest_error.NewBadRequestError("Invalid website url")
		}
	}

	return nil
}
