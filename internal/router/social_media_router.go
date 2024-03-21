package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zikri124/mygram-api/internal/handler"
	"github.com/zikri124/mygram-api/internal/middleware"
)

type SocialMediaRouter interface {
	Mount()
}

type socialMediaRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.SocialMediaHandler
}

func NewSocialMediaRouter(v *gin.RouterGroup, handler handler.SocialMediaHandler) SocialMediaRouter {
	return &socialMediaRouterImpl{v: v, handler: handler}
}

func (s *socialMediaRouterImpl) Mount() {
	s.v.Use(middleware.CheckAuth)
	s.v.POST("", s.handler.PostSocialMedia)
	s.v.GET("", s.handler.GetAllSocialMedias)
	s.v.PUT("/:id", s.handler.UpdateSocialMedia)
}
