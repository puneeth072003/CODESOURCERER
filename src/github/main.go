package main

import (
	"fmt"
	"github/controllers"

	"github/controllers/tokenhandlers"
	"github/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func initTokenManager() {
	envs, err := utils.Loadenv(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	installationID, exists := envs["INSTALLATION_ID"]
	if !exists {
		log.Fatalf("Error loading .env file: %v", err)
	}
	apiEndpoint := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID)

	jwtToken := tokenhandlers.GetJWT()
	tokenhandlers.NewTokenManager(apiEndpoint, jwtToken) // Initialize the TokenManager

	// Demo fetch of token
	token, err := tokenhandlers.GetInstance().GetToken()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Fetched Token:", token)
}

func main() {
	// Initialize the token manager
	initTokenManager()

	// Start the server
	router := gin.Default()
	router.GET("/code", controllers.Code)                // test route
	router.POST("/webhook", controllers.WebhookHandler)  // checking for push events
	router.GET("/testsend", controllers.TestSendPayload) // test route for payload generation
	router.Run(":3000")
}
