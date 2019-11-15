package ci_league

import (
	"context"
	"fmt"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
	"net/http"
	"sort"
	"time"
)

const owner = "mergermarket"

func GetIntegrations(repo string, githubToken string, idMappings map[string]string) TeamIntegrations {
	ctx := context.Background()
	client := github.NewClient(createOAauth2HTTPClient(githubToken))

	frequency := getCommitFrequency(client, ctx, repo, idMappings)
	integrations := sortedIntegrations(frequency)

	return integrations
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
		date = date.Add(-1 * (time.Hour *24))
	}

	year, month, day := date.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}

func getCommitFrequency(client *github.Client, ctx context.Context, repo string, idMappings map[string]string) map[Dev]int {
	options := github.CommitsListOptions{
		Since:       monday(),
		ListOptions: github.ListOptions{},
	}
	var allCommits []*github.RepositoryCommit
	for {
		commits, response, _ := client.Repositories.ListCommits(ctx, owner, repo, &options)
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
		} else {
			fmt.Println("couldnt find", name)
		}

		if name != "" {
			commitFrequency[Dev{
				Name:   name,
				Avatar: commit.GetAuthor().GetAvatarURL(),
			}]++
		}
	}

	return commitFrequency
}

func createOAauth2HTTPClient(githubToken string) *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return tc
}
