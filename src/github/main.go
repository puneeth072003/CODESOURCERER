package main

import (
	"fmt"
	"github/controllers"
	"github/controllers/finalizers"
	"github/controllers/tokenhandlers"
	"github/utils"
	"log"
	"time"

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
	var tokenManager = tokenhandlers.GetInstance()
	tokenManager = tokenhandlers.NewTokenManager(apiEndpoint, jwtToken)
	tokenManager.StartProactiveRefresh(1 * time.Minute) // Start proactive token refresh
}

func main() {
	// Initialize the token manager
	initTokenManager()

	// Demo fetch of token
	token, err := tokenhandlers.GetInstance().GetToken()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Fetched Token:", token)

	// Start the server
	router := gin.Default()
	router.GET("/code", controllers.Code)               // test route
	router.POST("/webhook", controllers.WebhookHandler) // checking for push events
	router.GET("/parse", finalizers.TestParseServer2Response)
	// router.GET("/testsend", controllers.TestSendPayload)
	router.Run(":3000")
}
