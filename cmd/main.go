package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zikri124/mygram-api/internal/handler"
	"github.com/zikri124/mygram-api/internal/infrastructure"
	"github.com/zikri124/mygram-api/internal/repository"
	"github.com/zikri124/mygram-api/internal/router"
	"github.com/zikri124/mygram-api/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Cannot load env: ", err)
	}

	g := gin.Default()

	gorm := infrastructure.NewGormPostgres()

	userRouteGroup := g.Group("/v1/users")
	userRepo := repository.NewUserRepository(gorm)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRouter := router.NewUserRouter(userRouteGroup, userHandler)
	userRouter.Mount()

	g.GET("/ping", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Server online"))
	})

	g.Run("127.0.0.1:3000")
}
