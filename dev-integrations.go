package ci_league

import "sort"

type Dev struct {
	Name   string
	Avatar string
}

type DevIntegrations struct {
	Integrations int
	Dev
}

type TeamIntegrations []DevIntegrations

func NewTeamIntegrations(frequencies map[Dev]int) TeamIntegrations {
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

func (t TeamIntegrations) Total() int {
	total := 0
	for _, integration := range t {
		total += integration.Integrations
	}
	return total
}
