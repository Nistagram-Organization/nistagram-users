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
	GetByUsername(*gin.Context)
	FollowUser(*gin.Context)
	CheckIfUserIsFollowing(*gin.Context)
	MuteUser(*gin.Context)
	CheckIfUserIsMuted(*gin.Context)
	BlockUser(*gin.Context)
	CheckIfUserIsBlocked(*gin.Context)
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

func (c *userController) GetByUsername(ctx *gin.Context) {
	email := ctx.Param("username")
	userEntity, err := c.usersService.GetByUsername(email)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}
	ctx.JSON(http.StatusOK, userEntity)
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

func (c *userController) FollowUser(ctx *gin.Context) {
	var followRequestDTO dtos.FollowRequestDTO
	if err := ctx.ShouldBindJSON(&followRequestDTO); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	followErr := c.usersService.FollowUser(&followRequestDTO)
	if followErr != nil {
		ctx.JSON(followErr.Status(), followErr)
		return
	}

	ctx.JSON(http.StatusOK, followErr)
}

func (c *userController) CheckIfUserIsFollowing(ctx *gin.Context) {
	following, checkErr := c.usersService.CheckIfUserIsFollowing(ctx.Query("user"), ctx.Query("following_user"))
	if checkErr != nil {
		ctx.JSON(checkErr.Status(), checkErr)
		return
	}

	ctx.JSON(http.StatusOK, following)
}

func (c *userController) MuteUser(ctx *gin.Context) {
	var muteDTO dtos.MuteDTO
	if err := ctx.ShouldBindJSON(&muteDTO); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	muteErr := c.usersService.MuteUser(&muteDTO)
	if muteErr != nil {
		ctx.JSON(muteErr.Status(), muteErr)
		return
	}

	ctx.JSON(http.StatusOK, muteErr)
}

func (c *userController) CheckIfUserIsMuted(ctx *gin.Context) {
	muted, checkErr := c.usersService.CheckIfUserIsMuted(ctx.Query("user"), ctx.Query("muted_user"))
	if checkErr != nil {
		ctx.JSON(checkErr.Status(), checkErr)
		return
	}

	ctx.JSON(http.StatusOK, muted)
}

func (c *userController) BlockUser(ctx *gin.Context) {
	var blockDTO dtos.BlockDTO
	if err := ctx.ShouldBindJSON(&blockDTO); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	blockErr := c.usersService.BlockUser(&blockDTO)
	if blockErr != nil {
		ctx.JSON(blockErr.Status(), blockErr)
		return
	}

	ctx.JSON(http.StatusOK, blockErr)
}

func (c *userController) CheckIfUserIsBlocked(ctx *gin.Context) {
	blocked, checkErr := c.usersService.CheckIfUserIsBlocked(ctx.Query("user"), ctx.Query("blocked_user"))
	if checkErr != nil {
		ctx.JSON(checkErr.Status(), checkErr)
		return
	}

	ctx.JSON(http.StatusOK, blocked)
}