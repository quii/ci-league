package ci_league_test

import (
	"context"
	ci_league "github.com/quii/ci-league"
	"github.com/quii/ci-league/github"
	"testing"
)

func TestGithubIntegrationsService_GetIntegrations(t *testing.T) {
	commitService := github.NewService(github.NewClient(""))
	service := ci_league.NewLeagueService(commitService, nil)
	_, err := service.GetStats(context.Background(), "quii", []string{"ci-league"})

	if err != nil {
		t.Fatalf("Failed to get integrations %s", err)
	}
}

func TestExtractAuthor(t *testing.T) {
	t.Run("findable", func(t *testing.T) {

		msg := `Remove default setting of anchor by useActiveNavItem

Co-authored-by: LisaMcCormack <lisamccormack85@gmail.com>`

		got := ci_league.ExtractCoAuthor(msg)
		want := "lisamccormack85@gmail.com"

		if got != want {
			t.Errorf("got %q, want %s", got, want)
		}
	})

	t.Run("cant find it", func(t *testing.T) {
		got := ci_league.ExtractCoAuthor("nope")
		want := ""

		if got != want {
			t.Errorf("got %q, want %s", got, want)
		}
	})
}
