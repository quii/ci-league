package ci_league

type Dev struct {
	Name   string
	Avatar string
}

type DevIntegrations struct {
	Integrations int
	Dev
}

type TeamIntegrations []DevIntegrations

func (t TeamIntegrations) Total() int {
	total := 0
	for _, integration := range t {
		total += integration.Integrations
	}
	return total
}

