package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func makeGitHubRequest(url, token string) ([]byte, error) {

	req, _ := http.NewRequest("GET", url, nil)
	configureJsonHeadersWithAuth(req, token)

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

func getJobID(jobUrl, token string) (int, error) {

	body, err := makeGitHubRequest(jobUrl, token)
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

func FetchLogs(jobUrl, owner, repo string) ([]string, error) {

	token, err := getRefreshToken()
	if err != nil {
		return nil, err
	}

	jobId, err := getJobID(jobUrl, token)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/jobs/%d/logs", owner, repo, jobId)
	body, err := makeGitHubRequest(url, token)
	if err != nil {
		return nil, err
	}

	logs := string(body)

	parsedLogs := make([]string, 1)
	for _, log := range strings.Split(logs, "\n") {
		if len(log) > 29 {
			parsedLogs = append(parsedLogs, log[29:]+"\n")
		}
	}

	return parsedLogs, nil
}
