package main

import (
	"encoding/json"
	"github/controllers"
	"github/controllers/finalizers"
	"github/utils"
	"io"
	"net/http"

	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func generateJWT(appID string, privkeyPath string) (string, error) {
	privkeyBytes, err := os.ReadFile(privkeyPath)
	if err != nil {
		return "", err
	}
	privkey, err := jwt.ParseRSAPrivateKeyFromPEM(privkeyBytes)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": appID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(privkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getInstallationAccessToken(installationID string, jwt string) (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to generate installation access token, status: %d, response: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return response.Token, nil
}

func main() {
	//load the envs
	envs, err := utils.Loadenv(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	privKeyPath, exists := envs["PRIVATE_KEY_PATH"]
	if !exists {
		log.Fatalf("Error loading .env file: %v", err)
	}
	appID, exists := envs["APP_ID"]
	if !exists {
		log.Fatalf("Error loading .env file: %v", err)
	}
	installationID, exists := envs["INSTALLATION_ID"]
	if !exists {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Get the private key path
	fmt.Println("PRIVATE_KEY_PATH:", privKeyPath)
	token, err := generateJWT(appID, string(privKeyPath))
	if err != nil {
		log.Fatalf("Error generating JWT: %v", err)
	}
	log.Printf("Generated a JWT: %s", token)

	installationAccessToken, err := getInstallationAccessToken(installationID, token)
	if err != nil {
		log.Fatalf("Error getting installation access token: %v", err)
	}
	log.Printf("\nGenerated an instalLation access token: %s", installationAccessToken)

	// Start the server
	router := gin.Default()
	router.GET("/code", controllers.Code)               // test route
	router.POST("/webhook", controllers.WebhookHandler) // checking for push events
	router.GET("/parse", finalizers.TestParseServer2Response)
	router.GET("/testsend", controllers.TestSendPayload)
	router.Run(":3000")
}
