package ci_league

import (
	"context"
	"time"
)

type SimpleCommit struct {
	Email     string
	AvatarURL string
	Message   string
	Status    string
}

type CommitService interface {
	GetCommits(ctx context.Context, since time.Time, owner string, repos ...string) ([]SimpleCommit, error)
}

type LeagueService struct {
	idMappings map[string]string
	commitService CommitService
}

func NewLeagueService(commitService CommitService, idMappings map[string]string) *LeagueService {
	return &LeagueService{idMappings: idMappings, commitService: commitService}
}

func (g *LeagueService) GetStats(ctx context.Context, owner string, repos []string) (TeamStats, error) {
	frequency, err := g.GetCommitFrequency(ctx, owner, repos)

	if err != nil {
		return nil, err
	}

	integrations := NewTeamStats(frequency)
	return integrations, nil
}

func (g *LeagueService) GetCommitFrequency(ctx context.Context, owner string, repos []string) (map[Dev]GitStat, error) {

	allCommits, err := g.commitService.GetCommits(ctx, monday(), owner, repos...)
	if err != nil {
		return nil, err
	}

	commitFrequency := make(map[string]int)
	failureFrequency := make(map[string]int)
	avatars := make(map[string]string)
	for _, commit := range allCommits {
		alias := g.findAlias(commit.Email)

		if commit.Status == "failure" {
			failureFrequency[alias]++
		}
		commitFrequency[alias]++
		avatars[alias] = commit.AvatarURL

		if coAuthor := ExtractCoAuthor(commit.Message); coAuthor != "" {
			coAuthor = g.findAlias(coAuthor)
			if commit.Status == "failure" {
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

func (g *LeagueService) findAlias(email string) string {
	if alias, found := g.idMappings[email]; found {
		return alias
	}
	return email
}

func monday() time.Time {
	date := time.Now()

	for date.Weekday() != time.Monday {
		date = date.Add(-1 * (time.Hour * 24))
	}

	year, month, day := date.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}
