package user

import (
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-users/src/services/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController interface {
	Edit(*gin.Context)
}

type userController struct {
	userService user.UserService
}

func NewUserController(userService user.UserService) UserController {
	return &userController{
		userService,
	}
}

func (c *userController) Edit(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	result, err := c.userService.Edit(&user)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
