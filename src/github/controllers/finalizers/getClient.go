package finalizers

import (
	"context"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

func GetClient(installationToken string) (*github.Client, context.Context) {
	ctx := context.Background()

	// Create an OAuth2 token source using the installation token
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: installationToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	return client, ctx
}
