package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v28/github"
	"github.com/quii/ci-league/league"
	"time"
)

type Service struct {
	client *github.Client
}

func NewService(client *github.Client) *Service {
	return &Service{client: client}
}

func (g *Service) GetCommits(ctx context.Context, since time.Time, owner string, repos ...string) ([]league.SimpleCommit, error) {
	var allCommits []league.SimpleCommit

	for _, repo := range repos {
		options := github.CommitsListOptions{
			Since:       since,
			ListOptions: github.ListOptions{},
		}
		for {
			commits, response, err := g.client.Repositories.ListCommits(ctx, owner, repo, &options)

			if err != nil {
				return nil, fmt.Errorf("couldn't get commits, %s", err)
			}

			for _, commit := range commits {
				simpleCommit, err := g.newSimpleCommit(ctx, owner, repo, commit)
				if err != nil {
					return allCommits, err
				}
				allCommits = append(allCommits, simpleCommit)
			}

			if response.NextPage == 0 {
				break
			}

			options.Page = response.NextPage
		}
	}
	return allCommits, nil
}

func (g *Service) newSimpleCommit(ctx context.Context, owner string, repo string, commit *github.RepositoryCommit) (league.SimpleCommit, error) {
	status, _, err := g.client.Repositories.GetCombinedStatus(ctx, owner, repo, commit.GetSHA(), nil)

	if err != nil {
		return league.SimpleCommit{}, fmt.Errorf("problem getting status %v", err)
	}

	simpleCommit := league.SimpleCommit{
		Email:     commit.GetCommit().GetAuthor().GetEmail(),
		AvatarURL: commit.GetAuthor().GetAvatarURL(),
		Message:   commit.GetCommit().GetMessage(),
		Status:    status.GetState(),
		CreatedAt: commit.GetCommit().GetAuthor().GetDate(),
	}

	return simpleCommit, nil
}
