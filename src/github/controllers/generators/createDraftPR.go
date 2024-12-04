package generators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github/utils"
	"io"
	"log"
	"net/http"
)

// Function to create a draft pull request
func CreateDraftPullRequest(repoOwner, repoName, branchName, commitSHA, title, body string) error {
	envs, err := utils.Loadenv(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	token, exists := envs["PAT_TOKEN"]
	if !exists {
		log.Fatalf("Error loading .env file: %v", err)
	}

	prURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", repoOwner, repoName)
	prData := map[string]interface{}{
		"title": title,
		"body":  body,
		"head":  branchName,
		"base":  "testing",
		"draft": true,
	}
	prPayload, err := json.Marshal(prData)
	if err != nil {
		return fmt.Errorf("unable to marshal PR data: %w", err)
	}

	req, err := http.NewRequest("POST", prURL, bytes.NewBuffer(prPayload))
	if err != nil {
		return fmt.Errorf("unable to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create draft PR, status: %s, response: %s", resp.Status, string(body))
	}

	log.Printf("Draft PR created successfully in repository %s/%s", repoOwner, repoName)
	return nil
}
