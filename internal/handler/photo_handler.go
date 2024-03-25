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

type PhotoHandler interface {
	PostPhoto(ctx *gin.Context)
	GetAllPhotosByUserId(ctx *gin.Context)
	GetPhotoById(ctx *gin.Context)
	UpdatePhoto(ctx *gin.Context)
	DeletePhoto(ctx *gin.Context)
}

type photoHandlerImpl struct {
	svc service.PhotoService
}

func NewPhotoHandler(svc service.PhotoService) PhotoHandler {
	return &photoHandlerImpl{svc: svc}
}

// Create Photo godoc
//
// @Summary		Create photo
// @Description	Create data to the login user
// @Tags		photo
// @Accept		json
// @Produce		json
// @Param		Authorization header 	string	true "Bearer token"
// @Param		photo	body		model.PhotoCreate	true	"New Photo"
// @Success		200		{object}	model.PhotoResCreate
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/photos [post]
func (p *photoHandlerImpl) PostPhoto(ctx *gin.Context) {
	userId, err := helper.GetUserIdFromGinCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	photoData := model.PhotoCreate{}
	err = ctx.ShouldBindJSON(&photoData)
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

// Get Photo godoc
//
// @Summary		Get all data of a photo by user id
// @Description	Return an array of photo data
// @Tags		photo
// @Accept		json
// @Produce		json
// @Param		Authorization header 	string	true "Bearer token"
// @Param       user_id    query    string  false  "user id of the owner"
// @Success		200		{object}	[]model.PhotoView
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/photos [get]
func (p *photoHandlerImpl) GetAllPhotosByUserId(ctx *gin.Context) {
	userIdStr := ctx.Request.URL.Query().Get("user_id")
	if userIdStr == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Missing User id in query"})
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	photos, err := p.svc.GetAllPhotosByUserId(ctx, uint32(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, photos)
}

// Get Photo godoc
//
// @Summary		Get data of a photo by photo id
// @Description	Get data by photo id
// @Tags		photo
// @Accept		json
// @Produce		json
// @Param		Authorization header string	true "Bearer token"
// @Param		id		path		int	true	"photo ID"
// @Success		200		{object}	model.PhotoView
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/photos/{id} [get]
func (p *photoHandlerImpl) GetPhotoById(ctx *gin.Context) {
	photoId, err := strconv.Atoi(ctx.Param("id"))
	if photoId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	photo, err := p.svc.GetPhotoById(ctx, uint32(photoId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, photo)
}

// Edit Photo godoc
//
// @Summary		Edit any photo data by photo id
// @Description	Edit any photo data by photo id
// @Tags		photo
// @Accept		json
// @Produce		json
// @Param		Authorization header string	true "Bearer token"
// @Param		id		path		int	true	"Photo id"
// @Param		photo	body		model.UpdatePhoto	true	"New Photo Editted"
// @Success		200		{object}	model.PhotoResUpdate
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/photos/{id} [put]
func (p *photoHandlerImpl) UpdatePhoto(ctx *gin.Context) {
	photoId, err := strconv.Atoi(ctx.Param("id"))
	if photoId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	photo, err := p.svc.GetPhotoById(ctx, uint32(photoId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Photo did not exist"})
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

	photoEditData := model.UpdatePhoto{}
	err = ctx.ShouldBindJSON(&photoEditData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(photoEditData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	photoUpdate := model.Photo{}
	photoUpdate.ID = uint32(photoId)
	photoUpdate.Title = photoEditData.Title
	photoUpdate.Caption = photoEditData.Caption
	photoUpdate.PhotoUrl = photoEditData.PhotoUrl

	photoRes, err := p.svc.UpdatePhoto(ctx, photoUpdate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, photoRes)
}

// Delete Photo godoc
//
// @Summary		Delete any photo
// @Description	Delete by id
// @Tags		photo
// @Accept		json
// @Produce		json
// @Param		Authorization header string	true "Bearer token"
// @Param		id		path		int	true	"Photo Id"
// @Success		200		{object}	response.SuccessResponse
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/photos/{id} [delete]
func (p *photoHandlerImpl) DeletePhoto(ctx *gin.Context) {
	photoId, err := strconv.Atoi(ctx.Param("id"))
	if photoId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	photo, err := p.svc.GetPhotoById(ctx, uint32(photoId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Photo did not exist"})
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

	err = p.svc.DeletePhoto(ctx, uint32(photoId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{Message: "Your photo has been successfully deleted"})
}
