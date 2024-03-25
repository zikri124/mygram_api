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
	auth    middleware.Authorization
}

func NewPhotoRouter(v *gin.RouterGroup, handler handler.PhotoHandler, auth middleware.Authorization) PhotoRouter {
	return &photoRouterImpl{v: v, handler: handler, auth: auth}
}

func (p *photoRouterImpl) Mount() {
	p.v.Use(p.auth.CheckAuth)
	p.v.POST("", p.handler.PostPhoto)
	p.v.GET("", p.handler.GetAllPhotosByUserId)
	p.v.GET("/:id", p.handler.GetPhotoById)
	p.v.PUT("/:id", p.handler.UpdatePhoto)
	p.v.DELETE("/:id", p.handler.DeletePhoto)
}
