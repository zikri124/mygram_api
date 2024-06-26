package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zikri124/mygram-api/internal/handler"
	"github.com/zikri124/mygram-api/internal/middleware"
)

type CommentRouter interface {
	Mount()
}

type commentRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.CommentHandler
	auth    middleware.Authorization
}

func NewCommentRouter(v *gin.RouterGroup, handler handler.CommentHandler, auth middleware.Authorization) CommentRouter {
	return &commentRouterImpl{v: v, handler: handler, auth: auth}
}

func (c *commentRouterImpl) Mount() {
	c.v.Use(c.auth.CheckAuth)
	c.v.POST("", c.handler.PostComment)
	c.v.GET("", c.handler.GetAllComments)
	c.v.GET("/:id", c.handler.GetCommentById)
	c.v.PUT("/:id", c.handler.UpdateComment)
	c.v.DELETE("/:id", c.handler.DeleteComment)
}
