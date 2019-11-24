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

type TeamStats struct {
	DevStats                             []DevStats
	TotalCommits, TotalFails, TotalScore int
}

func NewTeamStats(integrations map[Dev]GitStat) *TeamStats {
	team := &TeamStats{}
	for dev, stat := range integrations {
		team.DevStats = append(team.DevStats, DevStats{
			Dev:     dev,
			GitStat: stat,
		})

		team.TotalCommits += stat.Commits
		team.TotalFails += stat.Failures
		team.TotalScore += stat.Score()
	}

	sort.Slice(team.DevStats, func(i, j int) bool {
		if team.DevStats[i].Score() > team.DevStats[j].Score() {
			return true
		}
		if team.DevStats[i].Name > team.DevStats[i].Name {
			return true
		}
		return false
	})

	return team
}
