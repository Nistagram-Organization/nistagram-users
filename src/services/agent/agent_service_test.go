package agent

import (
	"errors"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/agent"
	"github.com/Nistagram-Organization/nistagram-users/src/repositories/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AgentServiceUnitTestsSuite struct {
	suite.Suite
	agentRepositoryMock *agent.AgentRepositoryMock
	userRepositoryMock  *user.UserRepositoryMock
	service             AgentService
}

func TestPostServiceUnitTestsSuite(t *testing.T) {
	suite.Run(t, new(AgentServiceUnitTestsSuite))
}

func (suite *AgentServiceUnitTestsSuite) SetupSuite() {
	suite.agentRepositoryMock = new(agent.AgentRepositoryMock)
	suite.userRepositoryMock = new(user.UserRepositoryMock)
	suite.service = NewAgentService(suite.agentRepositoryMock, suite.userRepositoryMock)
}

func (suite *AgentServiceUnitTestsSuite) TestNewAgentUserService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *AgentServiceUnitTestsSuite) TestAgentUserService_Create_RepositoryError() {
	agentEntity := model.Agent{}
	err := rest_error.NewInternalServerError("Error when trying to create agent", errors.New(""))

	suite.agentRepositoryMock.On("Create", &agentEntity).Return(nil, err).Once()

	_, userErr := suite.service.Create(&agentEntity)

	assert.Equal(suite.T(), err, userErr)
}

func (suite *AgentServiceUnitTestsSuite) TestAgentUserService_Create() {
	agentEntity := model.Agent{}

	suite.agentRepositoryMock.On("Create", &agentEntity).Return(&agentEntity, nil).Once()

	retUser, userErr := suite.service.Create(&agentEntity)

	assert.Equal(suite.T(), &agentEntity, retUser)
	assert.Equal(suite.T(), nil, userErr)
}
