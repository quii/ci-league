package main

import (
	"fmt"
	"github.com/google/go-github/v28/github"
	"github.com/quii/ci-league"
	"html/template"
	"log"
	"net/http"
	"os"
)

const port = ":8000"

func main() {

	var client *github.Client

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		client = github.NewClient(ci_league.NewOAauth2HTTPClient(token))
	} else {
		fmt.Println("Warning, providing no GITHUB_TOKEN env var means this will only work for public repos")
		client = github.NewClient(nil)
	}

	service := ci_league.NewGithubIntegrationsService(client, map[string]string{
		"tamara.jordan1+coding@hotmail.com":                         "Tamara",
		"27856297+dependabot-preview[bot]@users.noreply.github.com": "Depandabot",
		"qui666@gmail.com":            "Chris",
		"riyaddattani@gmail.com":      "Riya",
		"rick@22px.io":                "Ricky",
		"karol.slomczynski@gmail.com": "Osh",
		"riya_dattani@hotmail.com":    "Riya",
		"lisamccormack85@gmail.com":   "Lisa",
	})

	server := ci_league.NewServer(
		template.Must(template.ParseFiles("template.html")),
		service,
	)

	fmt.Println("Listening on port", port)
	if err := http.ListenAndServe(port, server); err != nil {
		log.Fatalf("Couldn't launch server listening on %s, %s", port, err)
	}
}
