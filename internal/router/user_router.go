package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zikri124/mygram-api/internal/handler"
	"github.com/zikri124/mygram-api/internal/middleware"
)

type UserRouter interface {
	Mount()
}

type userRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.UserHandler
	auth    middleware.Authorization
}

func NewUserRouter(v *gin.RouterGroup, handler handler.UserHandler, auth middleware.Authorization) UserRouter {
	return &userRouterImpl{v: v, handler: handler, auth: auth}
}

func (u *userRouterImpl) Mount() {
	u.v.GET("/:id", u.handler.GetUserById)
	u.v.POST("/register", u.handler.UserRegister)
	u.v.POST("/login", u.handler.UserLogin)
	u.v.Use(u.auth.CheckAuth)
	u.v.PUT("/:id", u.handler.UserEdit)
	u.v.DELETE("", u.handler.UserDelete)
}
