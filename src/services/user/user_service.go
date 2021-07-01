package user

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-users/src/dtos"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/post_user_repository"
	regUserRepo "github.com/Nistagram-Organization/nistagram-users/src/repositories/registered_user"
	user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
)

type UserService interface {
	AddPostToFavorites(*dtos.FavoritesDTO) rest_error.RestErr
	RemovePostFromFavorites(string, uint) rest_error.RestErr
	GetByEmail(string) (*user.User, rest_error.RestErr)
	Update(*user.User) (*user.User, rest_error.RestErr)
}

type userService struct {
	userRepository           user2.UserRepository
	registeredUserRepository regUserRepo.RegisteredUserRepository
	postUserRepository       post_user_repository.PostUserRepository
}

func NewUserService(userRepository user2.UserRepository, registeredUserRepository regUserRepo.RegisteredUserRepository,
	postUserRepository post_user_repository.PostUserRepository) UserService {
	return &userService{
		userRepository,
		registeredUserRepository,
		postUserRepository,
	}
}

func (s *userService) GetByEmail(email string) (*user.User, rest_error.RestErr) {
	return s.userRepository.GetByEmail(email)
}

func (s *userService) Update(editedUser *user.User) (*user.User, rest_error.RestErr) {
	existingUser, err := s.userRepository.GetByEmail(editedUser.Email)
	if err != nil {
		return nil, err
	}

	existingUser.Name = editedUser.Name
	existingUser.LastName = editedUser.LastName
	existingUser.Website = editedUser.Website
	existingUser.Phone = editedUser.Phone
	existingUser.Public = editedUser.Public
	existingUser.Taggable = editedUser.Taggable
	existingUser.Biography = editedUser.Biography

	if err := existingUser.Validate(); err != nil {
		return nil, err
	}

	updatedUser, err := s.userRepository.Update(existingUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) AddPostToFavorites(dto *dtos.FavoritesDTO) rest_error.RestErr {
	userEntity, userErr := s.userRepository.GetByEmail(dto.UserEmail)
	if userErr != nil {
		return userErr
	}

	for _, favorite := range userEntity.Favorites {
		if favorite.PostID == dto.PostID {
			return nil
		}
	}

	postUser := post.PostUser{
		PostID: dto.PostID,
		UserID: userEntity.ID,
	}

	userEntity.Favorites = append(userEntity.Favorites, postUser)

	_, err := s.userRepository.Update(userEntity)

	return err
}

func (s *userService) RemovePostFromFavorites(userMail string, postId uint) rest_error.RestErr {
	userEntity, userErr := s.userRepository.GetByEmail(userMail)
	if userErr != nil {
		return userErr
	}

	for _, favorite := range userEntity.Favorites {
		if favorite.PostID == postId {
			delErr := s.userRepository.DeleteFavorite(userEntity.ID, favorite.ID)
			if delErr != nil {
				return delErr
			}

			delErr = s.postUserRepository.Delete(favorite.ID)
			if delErr != nil {
				return delErr
			}
		}
	}

	_, err := s.userRepository.Update(userEntity)
	return err
}
