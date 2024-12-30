package finalizers

import (
	"context"
	"fmt"
)

func CreateFiles(client *github.Client, ctx context.Context, owner, repo, branch, filePath, content string) error {
	// Get the current file tree at the branch
	_, _, err := client.Repositories.CreateFile(ctx, owner, repo, filePath, github.RepositoryContentFileOptions{
		Committer: &github.CommitAuthor{
			Name:  github.String("Your Name"),
			Email: github.String("your-email@example.com"),
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
