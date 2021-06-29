package agent

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	agent2 "github.com/Nistagram-Organization/nistagram-users/src/repositories/agent"
	"net/mail"
	"net/url"
	"strings"
)

type AgentService interface {
	Create(agent *agent.Agent) (*agent.Agent, rest_error.RestErr)
	Delete(id uint) rest_error.RestErr
	Edit(agent *agent.Agent) (*agent.Agent, rest_error.RestErr)
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

func (s *agentService) Edit(agent *agent.Agent) (*agent.Agent, rest_error.RestErr) {
	if err := ValidateForEdit(agent); err != nil {
		return nil, err
	}

	existingAgent, err := s.agentRepository.GetByEmail(agent.Email)
	if err != nil {
		return nil, err
	}

	existingAgent.Name = agent.Name
	existingAgent.Surname = agent.Surname
	existingAgent.Website = agent.Website
	existingAgent.Phone = agent.Phone
	existingAgent.Public = agent.Public
	existingAgent.Taggable = agent.Taggable
	existingAgent.Biography = agent.Biography

	editedAgent, err := s.agentRepository.Edit(existingAgent)
	if err != nil {
		return nil, err
	}

	return editedAgent, nil
}

func ValidateForEdit(agent *agent.Agent) rest_error.RestErr {
	if _, err := mail.ParseAddress(agent.Email); err != nil {
		return rest_error.NewBadRequestError("Invalid email address")
	}
	if strings.TrimSpace(agent.Name) == "" {
		return rest_error.NewBadRequestError("Name cannot be empty")
	}
	if strings.TrimSpace(agent.Surname) == "" {
		return rest_error.NewBadRequestError("Surname cannot be empty")
	}
	if strings.TrimSpace(agent.Phone) == "" {
		return rest_error.NewBadRequestError("Phone cannot be empty")
	}
	if _, err := url.ParseRequestURI(agent.Website); err != nil {
		return rest_error.NewBadRequestError("Invalid website url")
	}

	return nil
}
