package finalizers

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v52/github"
)

func CreateFiles(client *github.Client, ctx context.Context, owner, repo, branch, filePath, content string) error {
	path, _ := strings.CutPrefix(filePath, "/")
	_, _, err := client.Repositories.CreateFile(ctx, owner, repo, path, &github.RepositoryContentFileOptions{ // Use & here
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
