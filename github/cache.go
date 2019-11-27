package github

import (
	"context"
	"fmt"
	"github.com/quii/ci-league/league"
	"io"
	"strings"
	"sync"
	"time"
)

type CachedService struct {
	delegate *Service
	out      io.Writer

	mutex sync.RWMutex
	cache map[commitRequest]*commitCache
}

func NewCachedService(delegate *Service, out io.Writer) *CachedService {
	return &CachedService{
		delegate: delegate,
		out:      out,
		cache:    make(map[commitRequest]*commitCache),
	}
}

type commitRequest struct {
	owner string
	repos string
}

func newCommitRequest(owner string, repos []string) commitRequest {
	return commitRequest{owner: owner, repos: strings.Join(repos, ",")}
}

type commitCache struct {
	lastUpdated time.Time
	commits     []league.SimpleCommit
}

func (c CachedService) GetCommits(ctx context.Context, since time.Time, owner string, repos ...string) ([]league.SimpleCommit, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	req := newCommitRequest(owner, repos)
	cache := c.cache[req]

	if cache == nil {
		cache = new(commitCache)
		c.cache[req] = cache
	}

	newSince := since
	if cache.lastUpdated.After(since) {
		newSince = cache.lastUpdated
		fmt.Fprintf(c.out, "Fetching since %s\n", newSince)
	} else {
		fmt.Fprintf(c.out, "Empty cache for %+v, fetching since %s\n", req, since)
	}

	newCommits, err := c.delegate.GetCommits(ctx, newSince, owner, repos...)

	if err != nil {
		return nil, err
	}

	cache.commits = append(cache.commits, newCommits...)
	cache.commits = removeOldCommits(cache.commits, since)
	cache.lastUpdated = time.Now()

	return cache.commits, nil
}

func removeOldCommits(commits []league.SimpleCommit, since time.Time) (commitsSince []league.SimpleCommit) {
	for _, commit := range commits {
		if commit.CreatedAt.After(since) {
			commitsSince = append(commitsSince, commit)
		}
	}
	return
}
