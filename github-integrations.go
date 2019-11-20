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

func (g *GithubIntegrationsService) GetIntegrations(ctx context.Context, owner string, repos []string) (TeamStats, error) {
	frequency, err := g.getCommitFrequency(ctx, owner, repos)

	if err != nil {
		return nil, err
	}

	integrations := NewTeamStats(frequency)
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

func (g *GithubIntegrationsService) findAlias(email string) string {
	if alias, found := g.idMappings[email]; found {
		return alias
	}
	return email
}

func (g *GithubIntegrationsService) getCommitFrequency(ctx context.Context, owner string, repos []string) (map[Dev]GitStat, error) {

	allCommits, err := g.getCommits(ctx, owner, repos...)
	if err != nil {
		return nil, err
	}

	commitFrequency := make(map[string]int)
	failureFrequency := make(map[string]int)
	avatars := make(map[string]string)
	for _, commit := range allCommits {
		alias := g.findAlias(commit.email)

		if commit.status == "failure" {
			failureFrequency[alias]++
		}
		commitFrequency[alias]++
		avatars[alias] = commit.avatarURL

		if coAuthor := ExtractCoAuthor(commit.message); coAuthor != "" {
			coAuthor = g.findAlias(coAuthor)
			if commit.status == "failure" {
				failureFrequency[coAuthor]++
			}
			commitFrequency[coAuthor]++
		}
	}

	devs := make(map[Dev]GitStat)

	for name, score := range commitFrequency {
		devs[Dev{
			Name:   name,
			Avatar: avatars[name],
		}] = GitStat{
			Commits:  score,
			Failures: failureFrequency[name],
		}
	}

	return devs, nil
}

type simpleCommit struct {
	email     string
	avatarURL string
	message   string
	status    string
}

func (g *GithubIntegrationsService) getCommits(ctx context.Context, owner string, repos ...string) ([]simpleCommit, error) {
	var allCommits []simpleCommit

	for _, repo := range repos {
		options := github.CommitsListOptions{
			Since:       monday(),
			ListOptions: github.ListOptions{},
		}
		for {
			commits, response, err := g.client.Repositories.ListCommits(ctx, owner, repo, &options)

			if err != nil {
				return nil, fmt.Errorf("couldn't get commits, %s", err)
			}

			for _, commit := range commits {
				status, _, err := g.client.Repositories.GetCombinedStatus(ctx, owner, repo, commit.GetSHA(), nil)

				if err != nil {
					return nil, fmt.Errorf("problem getting status %v", err)
				}
				allCommits = append(allCommits, simpleCommit{
					email:     commit.GetCommit().GetAuthor().GetEmail(),
					avatarURL: commit.GetAuthor().GetAvatarURL(),
					message:   commit.GetCommit().GetMessage(),
					status:    status.GetState(),
				})
			}

			if response.NextPage == 0 {
				break
			}

			options.Page = response.NextPage
		}
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
