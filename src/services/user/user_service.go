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
	GetById(uint) (*user.User, rest_error.RestErr)
	CheckIfPostIsInFavorites(string, uint) (bool, rest_error.RestErr)
	GetByUsername(string) (*user.User, rest_error.RestErr)
	CheckIfUserIsTaggable(string) bool
	FollowUser(*dtos.FollowRequestDTO) rest_error.RestErr
	CheckIfUserIsFollowing(string, string) (bool, rest_error.RestErr)
	GetFollowingUsers(string) ([]string, rest_error.RestErr)
	MuteUser(*dtos.MuteDTO) rest_error.RestErr
	CheckIfUserIsMuted(string, string) (bool, rest_error.RestErr)
	BlockUser(*dtos.BlockDTO) rest_error.RestErr
	CheckIfUserIsBlocked(string, string) (bool, rest_error.RestErr)
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

func (s *userService) GetByUsername(username string) (*user.User, rest_error.RestErr) {
	return s.userRepository.GetByUsername(username)
}

func (s *userService) Update(editedUser *user.User) (*user.User, rest_error.RestErr) {
	existingUser, err := s.userRepository.GetByEmail(editedUser.Email)
	if err != nil {
		return nil, err
	}

	existingUser.FirstName = editedUser.FirstName
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

	i := 0
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
		} else {
			userEntity.Favorites[i] = favorite
			i++
		}
	}
	userEntity.Favorites = userEntity.Favorites[:i]

	_, err := s.userRepository.Update(userEntity)
	return err
}

func (s *userService) CheckIfPostIsInFavorites(userEmail string, postID uint) (bool, rest_error.RestErr) {
	userEntity, userErr := s.userRepository.GetByEmail(userEmail)
	if userErr != nil {
		return false, userErr
	}

	for _, favorite := range userEntity.Favorites {
		if favorite.PostID == postID {
			return true, nil
		}
	}

	return false, nil
}

func (s *userService) CheckIfUserIsTaggable(username string) bool {
	userEntity, userErr := s.userRepository.GetByUsername(username)
	if userErr != nil {
		return false
	}

	return userEntity.Taggable
}

func (s *userService) FollowUser(followRequestDTO *dtos.FollowRequestDTO) rest_error.RestErr {
	userEntity, userErr := s.userRepository.GetByEmail(followRequestDTO.User)
	if userErr != nil {
		return userErr
	}

	userToFollow, userErr := s.userRepository.GetByEmail(followRequestDTO.UserToFollow)
	if userErr != nil {
		return userErr
	}

	for _, u := range userEntity.Following {
		if u.ID == userToFollow.ID {
			return nil
		}
	}

	userEntity.Following = append(userEntity.Following, *userToFollow)
	_, err := s.userRepository.Update(userEntity)

	return err
}

func (s *userService) CheckIfUserIsFollowing(userEmail string, userToFollowEmail string) (bool, rest_error.RestErr) {
	userEntity, userErr := s.userRepository.GetByEmail(userEmail)
	if userErr != nil {
		return false, userErr
	}

	userToFollow, userErr := s.userRepository.GetByEmail(userToFollowEmail)
	if userErr != nil {
		return false, userErr
	}

	for _, u := range userEntity.Following {
		if u.ID == userToFollow.ID {
			return true, nil
		}
	}

	return false, nil
}

func (s *userService) GetFollowingUsers(userEmail string) ([]string, rest_error.RestErr) {
	userEntity, userErr := s.userRepository.GetByEmail(userEmail)
	if userErr != nil {
		return nil, userErr
	}

	var followingUsers []string
	var muted bool
	var blocked bool
	for _, u := range userEntity.Following {
		muted, _ = s.CheckIfUserIsMuted(userEmail, u.Email)
		blocked, _ = s.CheckIfUserIsBlocked(userEmail, u.Email)

		if !muted && !blocked {
			followingUsers = append(followingUsers, u.Email)
		}
	}

	return followingUsers, nil
}

func (s *userService) GetById(id uint) (*user.User, rest_error.RestErr) {
	return s.userRepository.GetById(id)
}

func (s *userService) MuteUser(muteDTO *dtos.MuteDTO) rest_error.RestErr {
	userEntity, userErr := s.userRepository.GetByEmail(muteDTO.User)
	if userErr != nil {
		return userErr
	}

	userToMute, userErr := s.userRepository.GetByEmail(muteDTO.UserToMute)
	if userErr != nil {
		return userErr
	}

	for _, u := range userEntity.Muted {
		if u.ID == userToMute.ID {
			return nil
		}
	}

	userEntity.Muted = append(userEntity.Muted, *userToMute)
	_, err := s.userRepository.Update(userEntity)

	return err
}

func (s *userService) CheckIfUserIsMuted(userEmail string, mutedUser string) (bool, rest_error.RestErr) {
	userEntity, userErr := s.userRepository.GetByEmail(userEmail)
	if userErr != nil {
		return false, userErr
	}

	userToMute, userErr := s.userRepository.GetByEmail(mutedUser)
	if userErr != nil {
		return false, userErr
	}

	for _, u := range userEntity.Muted {
		if u.ID == userToMute.ID {
			return true, nil
		}
	}

	return false, nil
}

func (s *userService) BlockUser(blockDTO *dtos.BlockDTO) rest_error.RestErr {
	userEntity, userErr := s.userRepository.GetByEmail(blockDTO.User)
	if userErr != nil {
		return userErr
	}

	userToBlock, userErr := s.userRepository.GetByEmail(blockDTO.UserToBlock)
	if userErr != nil {
		return userErr
	}

	for _, u := range userEntity.Blocked {
		if u.ID == userToBlock.ID {
			return nil
		}
	}

	userEntity.Blocked = append(userEntity.Blocked, *userToBlock)
	_, err := s.userRepository.Update(userEntity)

	return err
}

func (s *userService) CheckIfUserIsBlocked(userEmail string, blockedUser string) (bool, rest_error.RestErr) {
	userEntity, userErr := s.userRepository.GetByEmail(userEmail)
	if userErr != nil {
		return false, userErr
	}

	userToBlock, userErr := s.userRepository.GetByEmail(blockedUser)
	if userErr != nil {
		return false, userErr
	}

	for _, u := range userEntity.Blocked {
		if u.ID == userToBlock.ID {
			return true, nil
		}
	}

	return false, nil
}