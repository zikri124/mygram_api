package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalln("Cannot load env: ", err)
	}

	g := gin.Default()

	g.GET("/ping", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Server online"))
	})

	g.Run("127.0.0.1:3000")
}
