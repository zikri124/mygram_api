package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/service"
	"github.com/zikri124/mygram-api/pkg/response"
)

type PhotoHandler interface {
	PostPhoto(ctx *gin.Context)
}

type photoHandlerImpl struct {
	svc service.PhotoService
}

func NewPhotoHandler(svc service.PhotoService) PhotoHandler {
	return &photoHandlerImpl{svc: svc}
}

func (p *photoHandlerImpl) PostPhoto(ctx *gin.Context) {
	userIdRaw, isExist := ctx.Get("UserId")
	if !isExist {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "payload not provided in access token"})
		return
	}

	userIdFloat := userIdRaw.(float64)
	userId := int(userIdFloat)
	if userId == 0 {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "payload not provided in access token"})
		return
	}

	photoData := model.PhotoCreate{}
	err := ctx.ShouldBindJSON(&photoData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(photoData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	photo := model.Photo{}
	photo.UserId = uint32(userId)
	photo.Caption = photoData.Caption
	photo.Title = photoData.Title
	photo.PhotoUrl = photoData.PhotoUrl

	photoRes, err := p.svc.PostPhoto(ctx, photo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, photoRes)
}
