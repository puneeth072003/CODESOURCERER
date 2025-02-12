package utils

import (
	"fmt"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomString(length int) string {
	rand.NewSource(time.Now().UnixNano())
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}

func GetRandomBranch() string {
	randomString := generateRandomString(5)
	return fmt.Sprintf("tests/CS-sandbox-%s", randomString)
}
