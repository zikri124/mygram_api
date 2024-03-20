package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zikri124/mygram-api/internal/handler"
	"github.com/zikri124/mygram-api/internal/middleware"
)

type PhotoRouter interface {
	Mount()
}

type photoRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.PhotoHandler
}

func NewPhotoRouter(v *gin.RouterGroup, handler handler.PhotoHandler) PhotoRouter {
	return &photoRouterImpl{v: v, handler: handler}
}

func (p *photoRouterImpl) Mount() {
	p.v.Use(middleware.CheckAuth)
	p.v.POST("", p.handler.PostPhoto)
	p.v.GET("", p.handler.GetAllPhotos)
}
