package finalizers

import (
	"context"

	"golang.org/x/oauth2"
)

func getClient(installationToken string) (*github.Client, context.Context) {
	ctx := context.Background()

	// Create an OAuth2 token source using the installation token
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: installationToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	return client, ctx
}
