package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v52/github"
)

func CreateBranch(client *github.Client, ctx context.Context, owner, repo, newBranchName string) error {
	// Get the default branch
	repoInfo, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return err
	}

	defaultBranch := repoInfo.GetDefaultBranch()

	// Get the latest commit on the default branch
	ref, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+defaultBranch)
	if err != nil {
		return err
	}

	// Create a new branch ref from the latest commit
	newRef := &github.Reference{
		Object: &github.GitObject{
			SHA: ref.Object.SHA,
		},
	}
	_, _, err = client.Git.CreateRef(ctx, owner, repo, &github.Reference{
		Ref:    github.String("refs/heads/" + newBranchName),
		Object: newRef.Object,
	})
	if err != nil {
		return err
	}

	fmt.Println("Created new branch:", newBranchName)
	return nil
}
