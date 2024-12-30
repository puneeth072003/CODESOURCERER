package initializers

import (
	"encoding/json"
	"fmt"
	"github/utils"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// FetchPullRequestDescription fetches the description of a pull request
func FetchPullRequestDescription(owner, repo string, prNumber int) (string, error) {
	owner, repo, err := utils.CleanURLParams(owner, repo, prNumber)
	if err != nil {
		return "", err
	}

	reqUrl, err := url.JoinPath("https://api.github.com", "repos", owner, repo, "pulls", strconv.Itoa(123))
	if err != nil {
		return "", fmt.Errorf("unable to construct request url: %v", err)
	}

	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

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

	var prData map[string]interface{}
	if err := json.Unmarshal(body, &prData); err != nil {
		return "", err
	}

	description, ok := prData["body"].(string)
	if !ok {
		return "", fmt.Errorf("unable to parse pull request description")
	}

	return description, nil
}
