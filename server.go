package main

import (
	"github.com/gin-gonic/gin"
	"webLimit/util"
)

func main() {
	server := gin.Default()
	//server.Use(util.GetHandler())
	server.Use(util.GetSlidingWindowHandler())
	server.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "root",
		})
	})
	server.Run(":8080")
}
