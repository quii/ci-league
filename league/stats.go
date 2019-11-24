package league

import (
	"sort"
)

type Dev struct {
	Name   string
	Avatar string
}

type GitStat struct {
	Commits  int
	Failures int
}

func (g *GitStat) Score() int {
	return g.Commits - (g.Failures * 3)
}

type DevStats struct {
	GitStat
	Dev
}

type TeamStats []DevStats

func NewTeamStats(integrations map[Dev]GitStat) TeamStats {
	var stats []DevStats
	for dev, stat := range integrations {
		stats = append(stats, DevStats{
			Dev:     dev,
			GitStat: stat,
		})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Score() > stats[j].Score()
	})
	return stats
}

func (t TeamStats) Total() (totalCommits int) {
	for _, integration := range t {
		totalCommits += integration.Commits
	}
	return
}

func (t TeamStats) TotalFails() (totalFails int) {
	for _, integration := range t {
		totalFails += integration.Failures
	}
	return
}

func (t TeamStats) TotalScore() (totalScore int) {
	for _, integration := range t {
		totalScore += integration.Score()
	}
	return
}
