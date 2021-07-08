package registered_user

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/registered_user"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-users/src/datasources/mysql"
	registered_user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/registered_user"
	user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type RegisteredUserServiceIntegrationTestsSuite struct {
	suite.Suite
	service RegisteredUserService
	db      *gorm.DB
}

func (suite *RegisteredUserServiceIntegrationTestsSuite) SetupSuite() {
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

	suite.service = NewRegisteredUserService(
		registeredUserRepository,
		userRepository,
	)
}

func TestRegisteredUserServiceIntegrationTestsSuite(t *testing.T) {
	suite.Run(t, new(RegisteredUserServiceIntegrationTestsSuite))
}

func (suite *RegisteredUserServiceIntegrationTestsSuite) TestIntegrationRegisteredUserService_Create() {
	userEntity := registered_user.RegisteredUser{}

	retUser, userErr := suite.service.Create(&userEntity)

	assert.Equal(suite.T(), &userEntity, retUser)
	assert.Equal(suite.T(), nil, userErr)
}
