package user

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/registered_user"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-users/src/datasources/mysql"
	"github.com/Nistagram-Organization/nistagram-users/src/dtos"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/post_user_repository"
	registered_user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/registered_user"
	user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type UserServiceIntegrationTestsSuite struct {
	suite.Suite
	service UserService
	db      *gorm.DB
	users   []model.User
}

func (suite *UserServiceIntegrationTestsSuite) SetupSuite() {
	database := mysql.NewMySqlDatabaseClient()
	if err := database.Init(); err != nil {
		panic(err)
	}

	if err := database.Migrate(
		&model.User{},
		&registered_user.RegisteredUser{},
		&agent.Agent{},
		&post.PostUser{},
	); err != nil {
		panic(err)
	}

	suite.db = database.GetClient()

	userRepository := user2.NewUserRepository(database)
	registeredUserRepository := registered_user2.NewRegisteredUserRepository(database)
	postUserRepository := post_user_repository.NewPostUserRepository(database)

	suite.service = NewUserService(
		userRepository,
		registeredUserRepository,
		postUserRepository,
	)
}

func (suite *UserServiceIntegrationTestsSuite) SetupTest() {
	suite.users = []model.User{
		{
			ID:       1,
			Username: "MAlen",
			Email:    "mejl@mail.com",
			Taggable: true,
		},
		{
			ID:       2,
			Username: "NAlen",
			Email:    "nedjo@mail.com",
			Taggable: true,
		},
	}
	postUserEntity := post.PostUser{
		ID:     1,
		PostID: 1,
		UserID: 1,
	}
	suite.users[0].Favorites = append(suite.users[0].Favorites, postUserEntity)

	tx := suite.db.Begin()
	tx.Create(&suite.users[0])
	tx.Create(&suite.users[1])
	tx.Commit()
}

func TestUserServiceIntegrationTestsSuite(t *testing.T) {
	suite.Run(t, new(UserServiceIntegrationTestsSuite))
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_GetByEmail_UserDoesNotExist() {
	email := "mail@mail.com"
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	_, userErr := suite.service.GetByEmail(email)

	assert.Equal(suite.T(), err, userErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_GetByEmail() {
	email := "mejl@mail.com"

	retUser, _ := suite.service.GetByEmail(email)

	assert.Equal(suite.T(), suite.users[0].ID, retUser.ID)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_GetByUsername_UserDoesNotExist() {
	username := "Juzernejm"
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	_, userErr := suite.service.GetByUsername(username)

	assert.Equal(suite.T(), err, userErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_GetByUsername() {
	username := "MAlen"

	retUser, _ := suite.service.GetByUsername(username)

	assert.Equal(suite.T(), suite.users[0], *retUser)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_AddPostToFavorites_PostDoesNotExist() {
	favoritesDTO := dtos.FavoritesDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	postErr := suite.service.AddPostToFavorites(&favoritesDTO)

	assert.Equal(suite.T(), err, postErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_AddPostToFavorites() {
	favoritesDTO := dtos.FavoritesDTO{
		PostID:    1,
		UserEmail: "mejl@mail.com",
	}

	postErr := suite.service.AddPostToFavorites(&favoritesDTO)

	assert.Equal(suite.T(), nil, postErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_RemovePostFromFavorites_UserDoesNotExist() {
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	postErr := suite.service.RemovePostFromFavorites("mail@mail.com", 1)

	assert.Equal(suite.T(), err, postErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_RemovePostFromFavorites() {
	postErr := suite.service.RemovePostFromFavorites("mejl@mail.com", 2)

	assert.Equal(suite.T(), nil, postErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_CheckIfPostIsInFavorites_UserDoesNotExist() {
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	_, postErr := suite.service.CheckIfPostIsInFavorites("mail@mail.com", 5)

	assert.Equal(suite.T(), err, postErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_CheckIfPostIsInFavorites() {
	inFavorites, postErr := suite.service.CheckIfPostIsInFavorites("mejl@mail.com", 5)

	assert.Equal(suite.T(), false, inFavorites)
	assert.Equal(suite.T(), nil, postErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_CheckIfUserIsTaggable_UserDoesNotExist() {
	isTaggable := suite.service.CheckIfUserIsTaggable("MujoAlen")

	assert.Equal(suite.T(), false, isTaggable)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_CheckIfUserIsTaggable() {
	isTaggable := suite.service.CheckIfUserIsTaggable("MAlen")

	assert.Equal(suite.T(), true, isTaggable)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_FollowUser_UserDoesNotExist() {
	followRequestDTO := dtos.FollowRequestDTO{
		User:         "mail@mail.com",
		UserToFollow: "nedjo@mail.com",
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	followErr := suite.service.FollowUser(&followRequestDTO)

	assert.Equal(suite.T(), err, followErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_FollowUser() {
	followRequestDTO := dtos.FollowRequestDTO{
		User:         "mejl@mail.com",
		UserToFollow: "nedjo@mail.com",
	}

	followErr := suite.service.FollowUser(&followRequestDTO)

	assert.Equal(suite.T(), nil, followErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_CheckIfUserIsFollowing_UserDoesNotExist() {
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	_, followErr := suite.service.CheckIfUserIsFollowing("mail@mail.com", "nedjo@mail.com")

	assert.Equal(suite.T(), err, followErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_CheckIfUserIsFollowing() {
	isFollowing, followErr := suite.service.CheckIfUserIsFollowing("nedjo@mail.com", "mejl@mail.com")

	assert.Equal(suite.T(), false, isFollowing)
	assert.Equal(suite.T(), nil, followErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_GetFollowingUsers_UserDoesNotExist() {
	err := rest_error.NewNotFoundError(fmt.Sprintf("User does not exist"))

	_, followErr := suite.service.GetFollowingUsers("mail@mail.com")

	assert.Equal(suite.T(), err, followErr)
}

func (suite *UserServiceIntegrationTestsSuite) TestIntegrationUserService_GetFollowingUsers() {
	followingUsers, followErr := suite.service.GetFollowingUsers("mejl@mail.com")

	assert.Equal(suite.T(), "nedjo@mail.com", followingUsers[0])
	assert.Equal(suite.T(), nil, followErr)
}