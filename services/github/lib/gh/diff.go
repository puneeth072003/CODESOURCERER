package gh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/codesourcerer-bot/github/utils"
)

// Fetch the list of changed files in the pull request
func FetchPullRequestFiles(owner, repo string, prNumber int) ([]map[string]interface{}, error) {
	owner, repo, err := utils.CleanURLParams(owner, repo, prNumber)
	if err != nil {
		return nil, err
	}

	reqUrl, err := url.JoinPath("https://api.github.com", "repos", owner, repo, "pulls", strconv.Itoa(prNumber), "files")
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
