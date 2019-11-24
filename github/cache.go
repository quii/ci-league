package github

import (
	"context"
	"fmt"
	"github.com/quii/ci-league/league"
	"io"
	"strings"
	"time"
)

type CachedService struct {
	delegate *Service
	cache    map[commitRequest]*commitCache
	out      io.Writer
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

type commitCache struct {
	lastUpdated time.Time
	commits     []league.SimpleCommit
}

func (c CachedService) GetCommits(ctx context.Context, since time.Time, owner string, repos ...string) ([]league.SimpleCommit, error) {
	req := commitRequest{
		owner: owner,
		repos: strings.Join(repos, ","),
	}
	cache := c.cache[req]

	if cache == nil {
		cache = new(commitCache)
		c.cache[req] = cache
	}

	if cache.lastUpdated.After(since) {
		since = cache.lastUpdated
		fmt.Fprintf(c.out, "Fetching since %s\n", since)
	} else {
		fmt.Fprintf(c.out, "Empty cache for %+v, fetching since %s\n", req, since)
	}

	newCommits, err := c.delegate.GetCommits(ctx, since, owner, repos...)

	if err != nil {
		return nil, err
	}

	cache.commits = append(cache.commits, newCommits...)
	cache.lastUpdated = time.Now()

	return cache.commits, nil
}
