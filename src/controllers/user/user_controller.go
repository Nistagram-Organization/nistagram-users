package user

import (
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
