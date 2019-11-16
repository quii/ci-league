package ci_league_test

import (
	"context"
	"github.com/google/go-github/v28/github"
	ci_league "github.com/quii/ci-league"
	"testing"
)

func TestGithubIntegrationsService_GetIntegrations(t *testing.T) {
	client := github.NewClient(nil)
	service := ci_league.NewGithubIntegrationsService(client, nil)
	_, err := service.GetIntegrations(context.Background(), "quii", "ci-league")

	if err != nil {
		t.Fatalf("Failed to get integrations %s", err)
	}
}
