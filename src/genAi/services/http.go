package services

import (
	"genAi/controllers"
	"genAi/models"
	"log"

	"github.com/gin-gonic/gin"
)

func StartGinServer(addr string) {

	ctx, client, model := models.InitializeModel()
	defer client.Close()

	// Initialize the Gin router
	router := gin.Default()

	router.GET("/ping", controllers.Pong)
	router.POST("/process", func(c *gin.Context) {
		controllers.ProcessAIRequest(ctx, c, client, model)
	})

	// Start the server on port 3001
	log.Println("Server running on port ", addr)
	router.Run(addr)

}
