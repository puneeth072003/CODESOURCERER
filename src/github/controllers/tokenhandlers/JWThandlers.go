package tokenhandlers

import (
	"github/utils"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GetJWT() string {
	//load envs
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

	// generate JWT
	token, err := GenerateJWT(appID, string(privKeyPath))
	if err != nil {
		log.Fatalf("Error generating JWT: %v", err)
	}
	return token
}

func GenerateJWT(appID string, privkeyPath string) (string, error) {
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
