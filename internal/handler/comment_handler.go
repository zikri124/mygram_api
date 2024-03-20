package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/zikri124/mygram-api/internal/model"
	"github.com/zikri124/mygram-api/internal/service"
	"github.com/zikri124/mygram-api/pkg/helper"
	"github.com/zikri124/mygram-api/pkg/response"
)

type CommentHandler interface {
	PostComment(ctx *gin.Context)
	GetAllComments(ctx *gin.Context)
}

type commentHandlerImpl struct {
	svc      service.CommentService
	photoSvc service.PhotoService
}

func NewCommentHandler(svc service.CommentService, photoSvc service.PhotoService) CommentHandler {
	return &commentHandlerImpl{svc: svc, photoSvc: photoSvc}
}

func (c *commentHandlerImpl) PostComment(ctx *gin.Context) {
	userId, err := helper.GetUserIdFromGinCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	newComment := model.CreateComment{}
	err = ctx.ShouldBindJSON(&newComment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(newComment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	photo, err := c.photoSvc.GetPhotoById(ctx, newComment.PhotoId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Photo did not exist"})
		return
	}

	comment := model.Comment{}
	comment.UserId = userId
	comment.Message = newComment.Message
	comment.PhotoId = newComment.PhotoId

	commentRes, err := c.svc.PostComment(ctx, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, commentRes)
}

func (c *commentHandlerImpl) GetAllComments(ctx *gin.Context) {
	comments, err := c.svc.GetAllComments(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
