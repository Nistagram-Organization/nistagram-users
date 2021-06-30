package user

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-users/src/dtos"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/post_user_repository"
	regUserRepo "github.com/Nistagram-Organization/nistagram-users/src/repositories/registered_user"
	user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
)

type UserService interface {
	AddPostToFavorites(*dtos.FavoritesDTO) rest_error.RestErr
	RemovePostFromFavorites(string, uint) rest_error.RestErr
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

	return s.userRepository.Update(userEntity)
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

	return s.userRepository.Update(userEntity)
}
