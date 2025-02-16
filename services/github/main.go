package main

import (
	"log"

	"github.com/codesourcerer-bot/github/controllers"
	"github.com/codesourcerer-bot/github/partials"
	"github.com/codesourcerer-bot/github/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	utils.LoadEnv()
	port := utils.GetPort()

	router := gin.Default()

	router.POST("/webhook", controllers.WebhookController)

	// Test Routes. Need to be removed later
	router.GET("/testsend", partials.TestSendPayload) // test route for payload generation
	router.GET("/testfinalizer", partials.TestFinalize)

	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
