package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
	"io"
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

func NewClient(token string, out io.Writer) *github.Client {
	var client *github.Client
	if token != "" {
		client = github.NewClient(NewOAauth2HTTPClient(token))
	} else {
		fmt.Fprintln(out, "Warning, providing no GITHUB_TOKEN env var means this will only work for public repos")
		client = github.NewClient(nil)
	}
	return client
}
