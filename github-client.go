package ci_league

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
)

func NewOAauth2HTTPClient(githubToken string) *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return tc
}
