package main

import (
	"log"

	"github.com/codesourcerer-bot/github/controllers"
	"github.com/codesourcerer-bot/github/lib/token"
	"github.com/codesourcerer-bot/github/partials"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the token manager
	token.InitTokenManager()

	router := gin.Default()

	router.POST("/webhook", controllers.WebhookController)

	// Test Routes. Need to be removed later
	router.GET("/testsend", partials.TestSendPayload) // test route for payload generation
	router.GET("/testfinalizer", partials.TestFinalize)

	if err := router.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
