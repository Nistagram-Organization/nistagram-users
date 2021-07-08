package registered_user

import (
	"errors"
	registered_user2 "github.com/Nistagram-Organization/nistagram-shared/src/model/registered_user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/registered_user"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RegisteredUserServiceUnitTestsSuite struct {
	suite.Suite
	registeredUserRepositoryMock *registered_user.RegisteredUserRepositoryMock
	userRepositoryMock           *user.UserRepositoryMock
	service                      RegisteredUserService
}

func TestPostServiceUnitTestsSuite(t *testing.T) {
	suite.Run(t, new(RegisteredUserServiceUnitTestsSuite))
}

func (suite *RegisteredUserServiceUnitTestsSuite) SetupSuite() {
	suite.userRepositoryMock = new(user.UserRepositoryMock)
	suite.registeredUserRepositoryMock = new(registered_user.RegisteredUserRepositoryMock)
	suite.service = NewRegisteredUserService(suite.registeredUserRepositoryMock, suite.userRepositoryMock)
}

func (suite *RegisteredUserServiceUnitTestsSuite) TestNewRegisteredUserService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *RegisteredUserServiceUnitTestsSuite) TestRegisteredUserService_Create_RepositoryError() {
	userEntity := registered_user2.RegisteredUser{}
	err := rest_error.NewInternalServerError("Error when trying to create registered user", errors.New(""))

	suite.registeredUserRepositoryMock.On("Create", &userEntity).Return(nil, err).Once()

	_, userErr := suite.service.Create(&userEntity)

	assert.Equal(suite.T(), err, userErr)
}

func (suite *RegisteredUserServiceUnitTestsSuite) TestRegisteredUserService_Create() {
	userEntity := registered_user2.RegisteredUser{}

	suite.registeredUserRepositoryMock.On("Create", &userEntity).Return(&userEntity, nil).Once()

	retUser, userErr := suite.service.Create(&userEntity)

	assert.Equal(suite.T(), &userEntity, retUser)
	assert.Equal(suite.T(), nil, userErr)
}
