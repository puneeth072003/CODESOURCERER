package token

import (
	"fmt"
	"log"
	"os"

	"github.com/codesourcerer-bot/github/lib"
)

func InitTokenManager() {

	installationID := os.Getenv("INSTALLATION_ID")
	if installationID == "" {
		log.Fatalf("no installation id found")
	}
	apiEndpoint := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID)

	jwtToken := lib.GetJWT()
	NewTokenManager(apiEndpoint, jwtToken) // Initialize the TokenManager

	// Demo fetch of token
	token, err := GetInstance().GetToken()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("Fetched Token:", token)
}
