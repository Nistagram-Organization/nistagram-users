package agent

import (
	"github.com/Nistagram-Organization/agent-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"gorm.io/gorm"
)

type AgentRepository interface {
	Create(*agent.Agent) (*agent.Agent, rest_error.RestErr)
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
