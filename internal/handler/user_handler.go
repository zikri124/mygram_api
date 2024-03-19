package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zikri124/mygram-api/internal/service"
)

type UserHandler interface {
	GetUserById(ctx *gin.Context)
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
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Failure		500		{object}	error
// @Router			/v1/users/{id} [get]
func (u *userHandlerImpl) GetUserById(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if userId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := u.svc.GetUserById(ctx, uint32(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	ctx.JSON(http.StatusOK, user)
}
