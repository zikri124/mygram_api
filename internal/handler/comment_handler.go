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
	GetCommentById(ctx *gin.Context)
	UpdateComment(ctx *gin.Context)
	DeleteComment(ctx *gin.Context)
}

type commentHandlerImpl struct {
	svc      service.CommentService
	photoSvc service.PhotoService
}

func NewCommentHandler(svc service.CommentService, photoSvc service.PhotoService) CommentHandler {
	return &commentHandlerImpl{svc: svc, photoSvc: photoSvc}
}

// Create Comment godoc
//
// @Summary		Create Comment
// @Description	Create data to the login user
// @Tags		comment
// @Accept		json
// @Produce		json
// @Param		Authorization header 	string	true "Bearer token"
// @Param		comment	body		model.CreateComment	true	"New Comment"
// @Success		200		{object}	model.CreateCommentRes
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/comments [post]
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

	commentRes, err := c.svc.PostComment(ctx, userId, newComment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, commentRes)
}

// Get Comment godoc
//
// @Summary		Get all data of a photo by user id
// @Description	Return an array of comments data
// @Tags		comment
// @Accept		json
// @Produce		json
// @Param		Authorization header 	string	true "Bearer token"
// @Param       photo_id    query    string  false  "photo id of the comments owner"
// @Success		200		{object}	[]model.CommentView
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/comments [get]
func (c *commentHandlerImpl) GetAllComments(ctx *gin.Context) {
	photoIdStr := ctx.Request.URL.Query().Get("photo_id")
	if photoIdStr == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Missing Photo id in query"})
		return
	}
	photoId, err := strconv.Atoi(photoIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	comments, err := c.svc.GetAllCommentsByPhotoId(ctx, uint32(photoId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

// Get Comment godoc
//
// @Summary		Get data of a comment by photo id
// @Description	Get data by comment id
// @Tags		comment
// @Accept		json
// @Produce		json
// @Param		Authorization header string	true "Bearer token"
// @Param		id		path		int	true	"Comment ID"
// @Success		200		{object}	model.CommentView
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/comments/{id} [get]
func (c *commentHandlerImpl) GetCommentById(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, comment)
}

// Edit Comment godoc
//
// @Summary		Edit any comment data by photo id
// @Description	Edit any comment data by photo id
// @Tags		comment
// @Accept		json
// @Produce		json
// @Param		Authorization header string	true "Bearer token"
// @Param		id		path		int	true	"Comment id"
// @Param		comment	body		model.UpdateComment	true	"New Comment Editted"
// @Success		200		{object}	model.UpdateCommentRes
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/comments/{id} [put]
func (c *commentHandlerImpl) UpdateComment(ctx *gin.Context) {
	commentId, err := strconv.Atoi(ctx.Param("id"))
	if commentId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	commentData, err := c.svc.GetCommentById(ctx, uint32(commentId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if commentData.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Comment did not exist"})
		return
	}

	userId, err := helper.GetUserIdFromGinCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if userId != uint32(commentData.UserId) {
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

	comment := model.Comment{ID: commentData.ID, UserId: commentData.UserId, PhotoId: commentData.PhotoId}
	comment.Message = commentEditData.Message

	commentRes, err := c.svc.UpdateComment(ctx, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, commentRes)
}

// Delete Comment godoc
//
// @Summary		Delete any comment
// @Description	Delete by id
// @Tags		comment
// @Accept		json
// @Produce		json
// @Param		Authorization header string	true "Bearer token"
// @Param		id		path		int	true	"Comment Id"
// @Success		200		{object}	response.SuccessResponse
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/comments/{id} [delete]
func (c *commentHandlerImpl) DeleteComment(ctx *gin.Context) {
	commentId, err := strconv.Atoi(ctx.Param("id"))
	if commentId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	photo, err := c.svc.GetCommentById(ctx, uint32(commentId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "comment did not exist"})
		return
	}

	userId, err := helper.GetUserIdFromGinCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if userId != uint32(photo.UserId) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "unauthorized to do this request"})
		return
	}

	err = c.svc.DeleteComment(ctx, uint32(commentId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{Message: "Your comment has been successfully deleted"})
}
