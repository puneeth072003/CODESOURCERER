package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Function to fetch the pull request comment from GitHub
func FetchPullRequestComment(repoOwner, repoName string, pullRequestNumber int) (string, error) {
	// Construct the URL for the PR comment API endpoint
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/comments", repoOwner, repoName, pullRequestNumber)

	// Send a GET request to the GitHub API to fetch the comments
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("unable to fetch pull request comments: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch PR comments, status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %v", err)
	}

	// Unmarshal the response into a slice of comments
	var comments []map[string]interface{}
	if err := json.Unmarshal(body, &comments); err != nil {
		return "", fmt.Errorf("unable to unmarshal PR comments: %v", err)
	}

	// Check if there are any comments available
	if len(comments) == 0 {
		return "", nil // No comments found
	}

	// Assuming the comment with dependencies is the last one (you can modify this based on your needs)
	// Return the body of the last comment
	return comments[len(comments)-1]["body"].(string), nil
}
