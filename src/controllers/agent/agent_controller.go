package agent

import (
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/agent"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-users/src/services/agent"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AgentController interface {
	Edit(*gin.Context)
}

type agentController struct {
	agentService agent.AgentService
}

func NewAgentController(agentService agent.AgentService) AgentController {
	return &agentController{
		agentService,
	}
}

func (c *agentController) Edit(ctx *gin.Context) {
	var user model.Agent
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	result, err := c.agentService.Edit(&user)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
