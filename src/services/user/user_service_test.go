package user

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-users/src/dtos"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/post_user_repository"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/registered_user"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserServiceUnitTestsSuite struct {
	suite.Suite
	userRepositoryMock           *user.UserRepositoryMock
	registeredUserRepositoryMock *registered_user.RegisteredUserRepositoryMock
	postUserRepositoryMock       *post_user_repository.PostUserRepositoryMock
	service                      UserService
}

func TestPostServiceUnitTestsSuite(t *testing.T) {
	suite.Run(t, new(UserServiceUnitTestsSuite))
}

func (suite *UserServiceUnitTestsSuite) SetupSuite() {
	suite.userRepositoryMock = new(user.UserRepositoryMock)
	suite.registeredUserRepositoryMock = new(registered_user.RegisteredUserRepositoryMock)
	suite.postUserRepositoryMock = new(post_user_repository.PostUserRepositoryMock)
	suite.service = NewUserService(suite.userRepositoryMock, suite.registeredUserRepositoryMock, suite.postUserRepositoryMock)
}

func (suite *UserServiceUnitTestsSuite) TestNewUserService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *UserServiceUnitTestsSuite) TestUserService_GetByEmail_UserDoesNotExist() {
	email := "mail@mail.com"
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	suite.userRepositoryMock.On("GetByEmail", email).Return(nil, err).Once()

	_, userErr := suite.service.GetByEmail(email)

	assert.Equal(suite.T(), err, userErr)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_GetByEmail() {
	email := "mail@mail.com"
	userEntity := model.User{}

	suite.userRepositoryMock.On("GetByEmail", email).Return(&userEntity, nil).Once()

	retUser, _ := suite.service.GetByEmail(email)

	assert.Equal(suite.T(), userEntity, *retUser)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_GetByUsername_UserDoesNotExist() {
	email := "mail@mail.com"
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	suite.userRepositoryMock.On("GetByUsername", email).Return(nil, err).Once()

	_, userErr := suite.service.GetByUsername(email)

	assert.Equal(suite.T(), err, userErr)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_GetByUsername() {
	email := "mail@mail.com"
	userEntity := model.User{}

	suite.userRepositoryMock.On("GetByUsername", email).Return(&userEntity, nil).Once()

	retUser, _ := suite.service.GetByUsername(email)

	assert.Equal(suite.T(), userEntity, *retUser)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_AddPostToFavorites_PostDoesNotExist() {
	favoritesDTO := dtos.FavoritesDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	suite.userRepositoryMock.On("GetByEmail", favoritesDTO.UserEmail).Return(nil, err).Once()

	postErr := suite.service.AddPostToFavorites(&favoritesDTO)

	assert.Equal(suite.T(), err, postErr)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_AddPostToFavorites_AlreadyAdded() {
	favoritesDTO := dtos.FavoritesDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}
	postUser := post.PostUser{
		PostID: 1,
	}
	userEntity := model.User{
		Favorites: []post.PostUser{postUser},
	}

	suite.userRepositoryMock.On("GetByEmail", favoritesDTO.UserEmail).Return(&userEntity, nil).Once()

	postErr := suite.service.AddPostToFavorites(&favoritesDTO)

	assert.Equal(suite.T(), nil, postErr)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_AddPostToFavorites() {
	favoritesDTO := dtos.FavoritesDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}
	postUser := post.PostUser{
		UserID: 1,
		PostID: 2,
	}
	userEntity := model.User{
		Favorites: []post.PostUser{postUser},
	}

	suite.userRepositoryMock.On("GetByEmail", favoritesDTO.UserEmail).Return(&userEntity, nil).Once()
	suite.userRepositoryMock.On("Update", &userEntity).Return(&userEntity, nil).Once()

	postErr := suite.service.AddPostToFavorites(&favoritesDTO)

	assert.Equal(suite.T(), nil, postErr)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_FollowUser_UserDoesNotExist() {
	followRequestDTO := dtos.FollowRequestDTO{
		User:         "mail@mail.com",
		UserToFollow: "mejl@mail.com",
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	suite.userRepositoryMock.On("GetByEmail", followRequestDTO.User).Return(nil, err).Once()

	followErr := suite.service.FollowUser(&followRequestDTO)

	assert.Equal(suite.T(), err, followErr)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_FollowUser_UserToFollowDoesNotExist() {
	followRequestDTO := dtos.FollowRequestDTO{
		User:         "mail@mail.com",
		UserToFollow: "mejl@mail.com",
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	suite.userRepositoryMock.On("GetByEmail", followRequestDTO.User).Return(&model.User{}, nil).Once()
	suite.userRepositoryMock.On("GetByEmail", followRequestDTO.UserToFollow).Return(nil, err).Once()

	followErr := suite.service.FollowUser(&followRequestDTO)

	assert.Equal(suite.T(), err, followErr)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_FollowUser_AlreadyFollowing() {
	followRequestDTO := dtos.FollowRequestDTO{
		User:         "mail@mail.com",
		UserToFollow: "mejl@mail.com",
	}
	userToFollow := model.User{
		ID: 1,
	}
	userEntity := model.User{
		Following: []model.User{userToFollow},
	}

	suite.userRepositoryMock.On("GetByEmail", followRequestDTO.User).Return(&userEntity, nil).Once()
	suite.userRepositoryMock.On("GetByEmail", followRequestDTO.UserToFollow).Return(&userToFollow, nil).Once()

	followErr := suite.service.FollowUser(&followRequestDTO)

	assert.Equal(suite.T(), nil, followErr)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_FollowUser() {
	followRequestDTO := dtos.FollowRequestDTO{
		User:         "mail@mail.com",
		UserToFollow: "mejl@mail.com",
	}
	userToFollow := model.User{
		ID: 1,
	}
	userEntity := model.User{
		Following: []model.User{},
	}

	suite.userRepositoryMock.On("GetByEmail", followRequestDTO.User).Return(&userEntity, nil).Once()
	suite.userRepositoryMock.On("GetByEmail", followRequestDTO.UserToFollow).Return(&userToFollow, nil).Once()
	suite.userRepositoryMock.On("Update", &userEntity).Return(&model.User{}, nil).Once()

	followErr := suite.service.FollowUser(&followRequestDTO)

	assert.Equal(suite.T(), nil, followErr)
}