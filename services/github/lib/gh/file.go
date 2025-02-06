package gh

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/v52/github"
)

func CreateFiles(client *github.Client, ctx context.Context, owner, repo, branch, filePath, content string) error {
	path, _ := strings.CutPrefix(filePath, "/")
	_, _, err := client.Repositories.CreateFile(ctx, owner, repo, path, &github.RepositoryContentFileOptions{ // Use & here
		Committer: &github.CommitAuthor{
			Name:  github.String("codesourcerer-bot"),
			Email: github.String("codesourcerer.org@gmail.com"),
		},
		Message: github.String("Adding new file " + filePath),
		Content: []byte(content),
		Branch:  github.String(branch),
	})
	if err != nil {
		return err
	}

	fmt.Println("File created:", filePath)
	return nil
}

func FetchFileFromGitHub(owner, repo, commitSHA, filePath string) (string, error) {
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
