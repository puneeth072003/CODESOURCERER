package utils

import (
	"crypto/rand"
	"fmt"
	"log"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = charset[int(bytes[i])%len(charset)]
	}

	return string(bytes), nil
}

func GetRandomBranch() string {
	randomString, err := generateRandomString(5)
	if err != nil {
		log.Fatalf("unable to generate random string: %v", err)
	}
	return fmt.Sprintf("tests/CS-sandbox-%s", randomString)
}
