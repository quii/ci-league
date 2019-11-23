package league_test

import (
	"context"
	"github.com/quii/ci-league/github"
	"github.com/quii/ci-league/league"
	"io/ioutil"
	"testing"
)

type stubAliasService string

func (s stubAliasService) GetAlias(string) string {
	return string(s)
}

func TestGithubIntegrationsService_GetIntegrations(t *testing.T) {
	commitService := github.NewService(github.NewClient("", ioutil.Discard))
	service := league.NewService(commitService, stubAliasService("Bob"))
	_, err := service.GetStats(context.Background(), "quii", []string{"ci-league"})

	if err != nil {
		t.Fatalf("Failed to get integrations %s", err)
	}
}
