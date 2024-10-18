package main

import (
	"github/controllers"
	"github/utils"

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
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(privkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
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
	// Get the private key path
	fmt.Println("PRIVATE_KEY_PATH:", privKeyPath)
	token, err := generateJWT("1028002", string(privKeyPath))
	if err != nil {
		log.Fatalf("Error generating JWT: %v", err)
	}
	log.Printf("Generated a JWT: %s", token)

	// Start the server
	router := gin.Default()
	router.GET("/ping", controllers.Pong)               // test route
	router.POST("/webhook", controllers.WebhookHandler) // checking for push events
	router.Run(":3000")
}
