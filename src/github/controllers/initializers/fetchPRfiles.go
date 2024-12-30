package initializers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

func cleanURLParams(owner, repo string, prNumber int) (string, string, error) {
	owner = url.QueryEscape(owner)
	repo = url.QueryEscape(repo)

	githubNameRegex := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

	if !githubNameRegex.MatchString(owner) || !githubNameRegex.MatchString(repo) || prNumber <= 0 {
		return "", "", fmt.Errorf("unable to clean url params")
	}

	return owner, repo, nil
}

// Fetch the list of changed files in the pull request
func FetchPullRequestFiles(owner, repo string, prNumber int) ([]map[string]interface{}, error) {

	owner, repo, err := cleanURLParams(owner, repo, prNumber)
	if err != nil {
		return nil, err
	}

	reqUrl, err := url.JoinPath("https://api.github.com", "repos", owner, repo, "pulls", strconv.Itoa(123), "files")
	if err != nil {
		return nil, fmt.Errorf("unable to construct request url: %v", err)
	}

	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API responded with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var files []map[string]interface{}
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, err
	}

	return files, nil
}
