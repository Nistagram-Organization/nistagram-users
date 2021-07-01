package user

import (
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-users/src/dtos"
	"github.com/Nistagram-Organization/nistagram-users/src/services/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserController interface {
	AddPostToFavorites(*gin.Context)
	RemovePostFromFavorites(*gin.Context)
	GetByEmail(*gin.Context)
	Update(*gin.Context)
}

type userController struct {
	usersService user.UserService
}

func NewUserController(usersService user.UserService) UserController {
	return &userController{
		usersService: usersService,
	}
}

func getId(idParam string) (uint, rest_error.RestErr) {
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, rest_error.NewBadRequestError("Id should be a number")
	}
	return uint(id), nil
}

func (c *userController) GetByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	user, err := c.usersService.GetByEmail(email)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *userController) Update(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	result, err := c.usersService.Update(&user)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *userController) AddPostToFavorites(ctx *gin.Context) {
	var favoritesDTO dtos.FavoritesDTO
	if err := ctx.ShouldBindJSON(&favoritesDTO); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	favErr := c.usersService.AddPostToFavorites(&favoritesDTO)
	if favErr != nil {
		ctx.JSON(favErr.Status(), favErr)
		return
	}

	ctx.JSON(http.StatusOK, favErr)
}

func (c *userController) RemovePostFromFavorites(ctx *gin.Context) {
	postId, idErr := getId(ctx.Query("post_id"))
	if idErr != nil {
		ctx.JSON(idErr.Status(), idErr)
		return
	}

	removeErr := c.usersService.RemovePostFromFavorites(ctx.Query("user_mail"), postId)
	if removeErr != nil {
		ctx.JSON(removeErr.Status(), removeErr)
		return
	}

	ctx.JSON(http.StatusOK, removeErr)
}
