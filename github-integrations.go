package ci_league

import (
	"context"
	"fmt"
	"github.com/google/go-github/v28/github"
	"sort"
	"time"
)

type GithubIntegrationsService struct {
	idMappings map[string]string
	client *github.Client
}

func NewGithubIntegrationsService(client *github.Client, idMappings map[string]string) *GithubIntegrationsService {
	return &GithubIntegrationsService{client: client, idMappings: idMappings}
}

func (g *GithubIntegrationsService) GetIntegrations(ctx context.Context, owner string, repo string) (TeamIntegrations, error) {
	frequency, err := getCommitFrequency(g.client, ctx, owner, repo, g.idMappings)

	if err != nil {
		return nil, err
	}

	integrations := sortedIntegrations(frequency)

	return integrations, nil
}

func sortedIntegrations(frequencies map[Dev]int) TeamIntegrations {
	var integrations []DevIntegrations
	for dev, count := range frequencies {
		integrations = append(integrations, DevIntegrations{
			Dev:          dev,
			Integrations: count,
		})
	}
	sort.Slice(integrations, func(i, j int) bool {
		return integrations[i].Integrations > integrations[j].Integrations
	})
	return integrations
}

func monday() time.Time {
	date := time.Now()

	for date.Weekday() != time.Monday {
		date = date.Add(-1 * (time.Hour * 24))
	}

	year, month, day := date.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}

func getCommitFrequency(client *github.Client, ctx context.Context, owner string, repo string, idMappings map[string]string) (map[Dev]int, error) {

	options := github.CommitsListOptions{
		Since:       monday(),
		ListOptions: github.ListOptions{},
	}
	var allCommits []*github.RepositoryCommit

	for {
		commits, response, err := client.Repositories.ListCommits(ctx, owner, repo, &options)

		if err != nil {
			return nil, fmt.Errorf("couldn't get commits, %s", err)
		}

		allCommits = append(allCommits, commits...)

		if response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage
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
