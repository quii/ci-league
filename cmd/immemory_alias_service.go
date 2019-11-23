package main

type InMemoryAliasService map[string]string

func (i InMemoryAliasService) GetAlias(email string) string {
	if alias, found := i[email]; found {
		return alias
	}
	return email
}
