package agent

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/mock"
)

type AgentRepositoryMock struct {
	mock.Mock
}

func (a *AgentRepositoryMock) Create(agentEntity *agent.Agent) (*agent.Agent, rest_error.RestErr) {
	args := a.Called(agentEntity)
	if args.Get(1) == nil {
		return args.Get(0).(*agent.Agent), nil
	}
	return nil, args.Get(1).(rest_error.RestErr)
}

func (a *AgentRepositoryMock) Delete(u uint) rest_error.RestErr {
	panic("implement me")
}

