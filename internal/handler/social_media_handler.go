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
