package handler

import (
	"net/http"
	"strconv"

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
	UpdateComment(ctx *gin.Context)
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

func (c *commentHandlerImpl) UpdateComment(ctx *gin.Context) {
	commentId, err := strconv.Atoi(ctx.Param("id"))
	if commentId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	comment, err := c.svc.GetCommentById(ctx, uint32(commentId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if comment.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Comment did not exist"})
		return
	}

	userId, err := helper.GetUserIdFromGinCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if userId != uint32(comment.UserId) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "unauthorized to do this request"})
		return
	}

	commentEditData := model.UpdateComment{}
	err = ctx.ShouldBindJSON(&commentEditData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(commentEditData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	comment.Message = commentEditData.Message

	commentRes, err := c.svc.UpdateComment(ctx, *comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, commentRes)
}
