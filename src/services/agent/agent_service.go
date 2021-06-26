package agent

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	agent2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/agent"
)

type AgentService interface {
	Create(agent *agent.Agent) (*agent.Agent, rest_error.RestErr)
	Delete(id uint) rest_error.RestErr
}

type agentService struct {
	agentRepository agent2.AgentRepository
}

func NewAgentService(agentRepository agent2.AgentRepository) AgentService {
	return &agentService{
		agentRepository,
	}
}

func (s *agentService) Create(agent *agent.Agent) (*agent.Agent, rest_error.RestErr) {
	return s.agentRepository.Create(agent)
}

func (s *agentService) Delete(id uint) rest_error.RestErr {
	return s.agentRepository.Delete(id)
}
