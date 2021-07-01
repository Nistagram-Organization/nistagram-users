package agent

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type AgentRepository interface {
	Create(*agent.Agent) (*agent.Agent, rest_error.RestErr)
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

func (r *agentRepository) Create(agent *agent.Agent) (*agent.Agent, rest_error.RestErr) {
	if err := r.db.Create(agent).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to create agent", err)
	}
	return agent, nil
}

func (r *agentRepository) Delete(id uint) rest_error.RestErr {
	if err := r.db.Delete(&agent.Agent{}, id).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to delete agent", err)
	}
	return nil
}
