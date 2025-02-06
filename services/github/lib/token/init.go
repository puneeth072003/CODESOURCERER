package token

import (
	"fmt"
	"log"

	"github.com/codesourcerer-bot/github/lib"
	"github.com/codesourcerer-bot/github/utils"
)

func InitTokenManager() {
	envs, err := utils.Loadenv(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	installationID, exists := envs["INSTALLATION_ID"]
	if !exists {
		log.Fatalf("Error loading .env file: %v", err)
	}
	apiEndpoint := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID)

	jwtToken := lib.GetJWT()
	NewTokenManager(apiEndpoint, jwtToken) // Initialize the TokenManager

	// Demo fetch of token
	token, err := GetInstance().GetToken()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Fetched Token:", token)
}
