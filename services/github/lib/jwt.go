package lib

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func getJWT() string {
	privKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	if privKeyPath == "" {
		log.Fatalf("private key path not found")
	}
	appID := os.Getenv("APP_ID")
	if appID == "" {
		log.Fatalf("app id not found")
	}

	// generate JWT
	token, err := generateJWT(appID, string(privKeyPath))
	if err != nil {
		log.Fatalf("Error generating JWT: %v", err)
	}
	return token
}

func generateJWT(appID string, privkeyPath string) (string, error) {
	privkeyBytes, err := os.ReadFile(filepath.Clean(privkeyPath))
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
		"exp": time.Now().Add(9 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(privkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
