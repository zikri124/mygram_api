package handler

import (
	"log"
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

	log.Println("hee")

	ctx.JSON(http.StatusOK, user)
}
