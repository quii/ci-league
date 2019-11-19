package ci_league

import (
	"context"
	"fmt"
	"github.com/google/go-github/v28/github"
	"regexp"
	"time"
)

type GithubIntegrationsService struct {
	idMappings map[string]string
	client     *github.Client
}

func NewGithubIntegrationsService(client *github.Client, idMappings map[string]string) *GithubIntegrationsService {
	return &GithubIntegrationsService{client: client, idMappings: idMappings}
}

func (g *GithubIntegrationsService) GetIntegrations(ctx context.Context, owner string, repo string) (TeamIntegrations, error) {
	frequency, err := g.getCommitFrequency(ctx, owner, repo, g.idMappings)

	if err != nil {
		return nil, err
	}

	integrations := NewTeamIntegrations(frequency)
	return integrations, nil
}

var (
	coAuthorRegex = regexp.MustCompile(`Co-authored-by:.*<(.*)>`)
)

func ExtractCoAuthor(message string) string {
	matches := coAuthorRegex.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (g *GithubIntegrationsService) getCommitFrequency(ctx context.Context, owner string, repo string, idMappings map[string]string) (map[Dev]int, error) {

	allCommits, err := g.getCommits(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	commitFrequency := make(map[Dev]int)
	for _, commit := range allCommits {

		name := commit.GetCommit().GetAuthor().GetEmail()

		if alias, found := idMappings[name]; found {
			name = alias
		}

		if name != "" {
			commitFrequency[Dev{
				Name:   name,
				Avatar: commit.GetAuthor().GetAvatarURL(),
			}]++
		}
	}

	return commitFrequency, nil
}

func (g *GithubIntegrationsService) getCommits(ctx context.Context, owner string, repo string) ([]*github.RepositoryCommit, error) {
	options := github.CommitsListOptions{
		Since:       monday(),
		ListOptions: github.ListOptions{},
	}
	var allCommits []*github.RepositoryCommit
	for {
		commits, response, err := g.client.Repositories.ListCommits(ctx, owner, repo, &options)

		if err != nil {
			return nil, fmt.Errorf("couldn't get commits, %s", err)
		}

		allCommits = append(allCommits, commits...)

		if response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage
	}
	return allCommits, nil
}

func monday() time.Time {
	date := time.Now()

	for date.Weekday() != time.Monday {
		date = date.Add(-1 * (time.Hour * 24))
	}

	year, month, day := date.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}
