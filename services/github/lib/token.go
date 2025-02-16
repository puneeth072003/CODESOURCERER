package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func getRefreshToken() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	installationID := os.Getenv("INSTALLATION_ID")
	if installationID == "" {
		log.Fatalf("no installation id found")
	}

	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	jwtToken := getJWT()
	configureJsonHeadersWithAuth(req, jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to refresh token: %v", resp.Status)
	}

	var response struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Token, nil
}
