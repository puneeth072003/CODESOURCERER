package finalizers

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func CreateFiles(client *github.Client, ctx context.Context, owner, repo, branch, filePath, content string) error {
	// Get the current file tree at the branch
	_, _, err := client.Repositories.CreateFile(ctx, owner, repo, filePath, &github.RepositoryContentFileOptions{ // Use & here
		Committer: &github.CommitAuthor{
			Name:  github.String("CODESOURCERER"),
			Email: github.String("pyd773@gmail.com"),
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
