package main

import (
	"genAi/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", controllers.Pong)
	router.Run(":3001")
}
