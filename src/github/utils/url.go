package utils

import (
	"fmt"
	"net/url"
	"regexp"
)

func CleanURLParams(owner, repo string, prNumber int) (string, string, error) {
	owner = url.QueryEscape(owner)
	repo = url.QueryEscape(repo)

	githubNameRegex := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

	if !githubNameRegex.MatchString(owner) || !githubNameRegex.MatchString(repo) || prNumber <= 0 {
		return "", "", fmt.Errorf("unable to clean url params")
	}

	return owner, repo, nil
}
