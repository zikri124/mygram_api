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

type SocialMediaHandler interface {
	PostSocialMedia(ctx *gin.Context)
	GetAllSocialMediasByUserId(ctx *gin.Context)
	GetSocialMediaById(ctx *gin.Context)
	UpdateSocialMedia(ctx *gin.Context)
	DeleteSocialMedia(ctx *gin.Context)
}

type socialMediaHandlerImpl struct {
	svc service.SocialMediaService
}

func NewSocialMediaHandler(svc service.SocialMediaService) SocialMediaHandler {
	return &socialMediaHandlerImpl{svc: svc}
}

// Create Social Media godoc
//
// @Summary		Create social_media
// @Description	Create data to the login user
// @Tags		social_media
// @Accept		json
// @Produce		json
// @Param		Authorization header 	string	true "Bearer token"
// @Param		photo	body		model.NewSocialMedia	true	"New Social Media"
// @Success		200		{object}	model.CreateSocialMediaRes
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/social_medias [post]
func (s *socialMediaHandlerImpl) PostSocialMedia(ctx *gin.Context) {
	userId, err := helper.GetUserIdFromGinCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	newSocial := model.NewSocialMedia{}
	err = ctx.ShouldBindJSON(&newSocial)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(newSocial)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	socialMediaRes, err := s.svc.PostSocial(ctx, userId, newSocial)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, socialMediaRes)
}

// Get Social Media godoc
//
// @Summary		Get all data of a social_media by user id
// @Description	Return an array of social_media data
// @Tags		social_media
// @Accept		json
// @Produce		json
// @Param		Authorization header 	string	true "Bearer token"
// @Param       user_id    query    string  false  "user id of the owner"
// @Success		200		{object}	[]model.SocialMediaView
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/social_medias [get]
func (s *socialMediaHandlerImpl) GetAllSocialMediasByUserId(ctx *gin.Context) {
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

	socials, err := s.svc.GetAllSocialMediasByUserId(ctx, uint32(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, socials)
}

// Get Social Media godoc
//
// @Summary		Get data of a social_media by social_media id
// @Description	Get data by social_media id
// @Tags		social_media
// @Accept		json
// @Produce		json
// @Param		Authorization header string	true "Bearer token"
// @Param		id		path		int	true	"Social Media ID"
// @Success		200		{object}	model.SocialMediaView
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/social_medias/{id} [get]
func (s *socialMediaHandlerImpl) GetSocialMediaById(ctx *gin.Context) {
	socialId, err := strconv.Atoi(ctx.Param("id"))
	if socialId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	social, err := s.svc.GetSocialById(ctx, uint32(socialId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if social.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "User Social Media data did not exist"})
		return
	}

	ctx.JSON(http.StatusOK, social)
}

// Edit Social Media godoc
//
// @Summary		Edit any social_media data by social_media id
// @Description	Edit any social_media data by social_media id
// @Tags		social_media
// @Accept		json
// @Produce		json
// @Param		Authorization header string	true "Bearer token"
// @Param		id		path		int	true	"Photo id"
// @Param		social_media	body		model.NewSocialMedia	true	"New Social Media Editted"
// @Success		200		{object}	model.UpdateSocialMediaRes
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/social_medias/{id} [put]
func (s *socialMediaHandlerImpl) UpdateSocialMedia(ctx *gin.Context) {
	socialId, err := strconv.Atoi(ctx.Param("id"))
	if socialId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	socialData, err := s.svc.GetSocialById(ctx, uint32(socialId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if socialData.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "User Social Media data did not exist"})
		return
	}

	userId, err := helper.GetUserIdFromGinCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if userId != uint32(socialData.UserId) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "unauthorized to do this request"})
		return
	}

	socialUpdateData := model.NewSocialMedia{}
	err = ctx.ShouldBindJSON(&socialUpdateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(socialUpdateData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	social := model.SocialMedia{ID: socialData.ID, UserId: socialData.UserId}
	social.Name = socialUpdateData.Name
	social.SocialMediaUrl = socialUpdateData.SocialMediaUrl

	socialMediaRes, err := s.svc.UpdateSocial(ctx, social)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, socialMediaRes)
}

// Delete Social Media godoc
//
// @Summary		Delete any social_media
// @Description	Delete by social_media id
// @Tags		social_media
// @Accept		json
// @Produce		json
// @Param		Authorization header string	true "Bearer token"
// @Param		id		path		int	true	"Social Media Id"
// @Success		200		{object}	response.SuccessResponse
// @Failure		400		{object}	response.ErrorResponse
// @Failure		500		{object}	response.ErrorResponse
// @Router		/v1/social_medias/{id} [delete]
func (s *socialMediaHandlerImpl) DeleteSocialMedia(ctx *gin.Context) {
	socialId, err := strconv.Atoi(ctx.Param("id"))
	if socialId == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	social, err := s.svc.GetSocialById(ctx, uint32(socialId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if social.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "User Social Media data did not exist"})
		return
	}

	userId, err := helper.GetUserIdFromGinCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if userId != uint32(social.UserId) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "unauthorized to do this request"})
		return
	}

	err = s.svc.DeleteSocial(ctx, uint32(socialId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{Message: "Your social media has been successfully deleted"})
}
