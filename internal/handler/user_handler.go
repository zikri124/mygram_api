package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/service"
	"github.com/zikri124/mygram-api/pkg/response"
)

type UserHandler interface {
	GetUserById(ctx *gin.Context)
	UserRegister(ctx *gin.Context)
	UserLogin(ctx *gin.Context)
}

type userHandlerImpl struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return &userHandlerImpl{svc: svc}
}

// Show User by Id godoc
//
// @Summary		Show user data by user id
// @Description	Show data of user by id given in params
// @Tags			users
// @Accept			json
// @Produce		json
// @Param		id		path		int		true	"User ID"
// @Success		200		{object}	model.UserView
// @Failure		400		{object}	response.ErrorResponse
// @Failure		404		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router			/v1/users/{id} [get]
func (u *userHandlerImpl) GetUserById(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := u.svc.GetUserById(ctx, uint32(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Register User godoc
//
// @Summary		Register a new user
// @Description	Register a new user to
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user	body	model.UserSignUp	true	"New User"
// @Success		200		{object}	model.UserView
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/users/register [post]
func (u *userHandlerImpl) UserRegister(ctx *gin.Context) {
	userRegData := model.UserSignUp{}

	err := ctx.ShouldBindJSON(&userRegData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(userRegData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := u.svc.UserRegister(ctx, userRegData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// Login User godoc
//
// @Summary		Route to login user
// @Description	If success, login route return an access token
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user	body	model.UserSignIn	true	"Login User"
// @Success		200		{object}	response.TokenResponse
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/users/login [post]
func (u *userHandlerImpl) UserLogin(ctx *gin.Context) {
	userData := model.UserSignIn{}
	err := ctx.ShouldBindJSON(&userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(userData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := u.svc.UserLogin(ctx, userData)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: err.Error()})
		return
	}

	token, err := u.svc.GenerateAccessToken(ctx, *user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response.TokenResponse{Token: token})
}
