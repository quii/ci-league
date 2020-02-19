package league

import (
	"regexp"
)

var (
	coAuthorRegex = regexp.MustCompile(`<(.*?)\>`)
)

func extractCoAuthor(message string) []string {
	matches := coAuthorRegex.FindAllStringSubmatch(message, -1)
	var emails []string
	if len(matches) > 0 {
		for _, email := range matches {
			emails = append(emails, email[1])
		}
		return emails
	}
	return []string{}
}
