package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

	g.GET("/ping", func(ctx *gin.Context) {
		log.Println("Server online")
		ctx.Writer.Write([]byte("Server online"))
	})

	g.Run("127.0.0.1:3000")
}
