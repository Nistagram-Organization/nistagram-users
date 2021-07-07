package agent

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/registered_user"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-users/src/datasources/mysql"
	agent2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/agent"
	user2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type AgentServiceIntegrationTestsSuite struct {
	suite.Suite
	service AgentService
	db      *gorm.DB
}

func (suite *AgentServiceIntegrationTestsSuite) SetupSuite() {
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
	agentRepository := agent2.NewAgentRepository(database)

	suite.service = NewAgentService(
		agentRepository,
		userRepository,
	)
}

func TestAgentServiceIntegrationTestsSuite(t *testing.T) {
	suite.Run(t, new(AgentServiceIntegrationTestsSuite))
}

func (suite *AgentServiceIntegrationTestsSuite) TestIntegrationAgentService_Create() {
	userEntity := agent.Agent{}

	retUser, userErr := suite.service.Create(&userEntity)

	assert.Equal(suite.T(), &userEntity, retUser)
	assert.Equal(suite.T(), nil, userErr)
}
