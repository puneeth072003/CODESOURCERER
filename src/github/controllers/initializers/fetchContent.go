package initializers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func FetchFileContentFromGitHub(owner, repo, commitSHA, filePath string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, filePath, commitSHA)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	// creates a pointer to new http client instance (struct)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API responded with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf("Successfully fetched content for file: %s", filePath)
	return string(body), nil
}
