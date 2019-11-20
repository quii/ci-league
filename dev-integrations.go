package ci_league

import "sort"

type Dev struct {
	Name   string
	Avatar string
}

type GitStat struct {
	Commits int
	Failures int
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
			Dev:          dev,
			GitStat: stat,
		})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Commits > stats[j].Commits
	})
	return stats
}

func (t TeamStats) Total() int {
	total := 0
	for _, integration := range t {
		total += integration.Commits
	}
	return total
}
