package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zikri124/mygram-api/pkg/helper"
	"github.com/zikri124/mygram-api/pkg/response"
)

func CheckAuth(ctx *gin.Context) {
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
	ctx.Next()
}
