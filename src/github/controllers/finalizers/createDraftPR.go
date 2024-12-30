package finalizers

import (
	"context"
	"fmt"

	"github.com/google/go-github/v52/github"
)

func CreateDraftPR(client *github.Client, ctx context.Context, owner, repo, title, headBranch, baseBranch, body string) error {
	pr := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(headBranch),
		Base:  github.String(baseBranch),
		Body:  github.String(body),
		Draft: github.Bool(true),
	}

	_, _, err := client.PullRequests.Create(ctx, owner, repo, pr)
	if err != nil {
		return err
	}

	fmt.Println("Draft Pull Request created:", title)
	return nil
}
