package gh

import (
	"context"
	"log"

	"github.com/google/go-github/v52/github"
)

func CreateBranch(client *github.Client, ctx context.Context, owner, repo, baseBranch, newBranchName string) error {
	// Get the latest commit on the default branch
	ref, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+baseBranch)
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

	log.Println("Created new branch:", newBranchName)
	return nil
}
