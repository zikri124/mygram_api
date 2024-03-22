package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zikri124/mygram-api/internal/service"
	"github.com/zikri124/mygram-api/pkg/helper"
	"github.com/zikri124/mygram-api/pkg/response"
)

type Authorization interface {
	CheckAuth(ctx *gin.Context)
}

type authorizationImpl struct {
	userService service.UserService
}

func NewAuthorization(userService service.UserService) Authorization {
	return &authorizationImpl{userService: userService}
}

func (a *authorizationImpl) CheckAuth(ctx *gin.Context) {
	auth := ctx.GetHeader("Authorization")

	authArr := strings.Split(auth, " ")
	if len(authArr) < 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Message: "unauthorized",
			Errors: []string{"invalid token"},
		})
		return
	}

	if authArr[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Message: "unauthorized",
			Errors: []string{"invalid token"},
		})
		return
	}

	token := authArr[1]
	claims, err := helper.ValidateToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Message: "unauthorized",
			Errors: []string{"invalid token", "failed to decode token"},
		})
		return
	}

	ctx.Set("UserId", claims["user_id"])

	userId, err := helper.GetUserIdFromGinCtx(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse{Message: "error when get user id from token"})
		return
	}

	user, err := a.userService.GetUserById(ctx, uint32(userId))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse{Message: "error when get user id from token"})
		return
	}

	if user.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Message: "unauthorized",
			Errors: []string{"invalid token"},
		})
		return
	}

	ctx.Next()
}
