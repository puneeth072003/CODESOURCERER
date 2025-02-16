package lib

import (
	"context"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

func GetClient() (*github.Client, context.Context, error) {
	ctx := context.Background()

	refreshToken, err := getRefreshToken()
	if err != nil {
		return nil, nil, err
	}

	// Create an OAuth2 token source using the installation token
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: refreshToken})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	return client, ctx, nil
}
