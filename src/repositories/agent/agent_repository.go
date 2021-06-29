package agent

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type AgentRepository interface {
	GetByEmail(email string) (*agent.Agent, rest_error.RestErr)
	Create(*agent.Agent) (*agent.Agent, rest_error.RestErr)
	Edit(*agent.Agent) (*agent.Agent, rest_error.RestErr)
	Delete(uint) rest_error.RestErr
}

type agentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(databaseClient datasources.DatabaseClient) AgentRepository {
	return &agentRepository{
		databaseClient.GetClient(),
	}
}

func (r *agentRepository) GetByEmail(email string) (*agent.Agent, rest_error.RestErr) {
	agent := agent.Agent{
		User: user.User{
			Email: email,
		},
	}
	if err := r.db.Take(&agent, agent.User.Email).Error; err != nil {
		return nil, rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get agent with email %s", agent.User.Email))
	}
	return &agent, nil
}

func (r *agentRepository) Create(agent *agent.Agent) (*agent.Agent, rest_error.RestErr) {
	if err := r.db.Create(agent).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to create agent", err)
	}
	return agent, nil
}

func (r *agentRepository) Edit(agent *agent.Agent) (*agent.Agent, rest_error.RestErr) {
	if err := r.db.Save(agent).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to edit agent", err)
	}
	return agent, nil
}

func (r *agentRepository) Delete(id uint) rest_error.RestErr {
	if err := r.db.Delete(&agent.Agent{}, id).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to delete agent", err)
	}
	return nil
}
