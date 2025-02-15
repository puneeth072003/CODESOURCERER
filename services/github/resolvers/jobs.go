package resolvers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/codesourcerer-bot/github/lib/token"
)

func makeGitHubRequest(url string) ([]byte, error) {
	token, err := token.GetInstance().GetToken()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API request failed: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func getJobID(jobUrl string) (int, error) {

	body, err := makeGitHubRequest(jobUrl)
	if err != nil {
		return 0, err
	}

	var result struct {
		Jobs []struct {
			ID int `json:"id"`
		} `json:"jobs"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil || len(result.Jobs) == 0 {
		return 0, fmt.Errorf("no jobs found: %v", err)
	}

	return result.Jobs[0].ID, nil
}

func FetchLogs(jobUrl, owner, repo string) (string, error) {

	jobId, err := getJobID(jobUrl)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/jobs/%d/logs", owner, repo, jobId)
	body, err := makeGitHubRequest(url)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
