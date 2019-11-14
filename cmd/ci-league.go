package main

import (
	"github.com/quii/ci-league"
	"html/template"
	"log"
	"net/http"
	"os"
)

const repo = "deals-page-subscriber"

var (
	githubEmailMappings = map[string]string{
		"tamara.jordan1+coding@hotmail.com":                         "Tamara",
		"27856297+dependabot-preview[bot]@users.noreply.github.com": "Depandabot",
		"qui666@gmail.com":            "Chris",
		"riyaddattani@gmail.com":      "Riya",
		"rick@22px.io":                "Ricky",
		"karol.slomczynski@gmail.com": "Osh",
		"riya_dattani@hotmail.com":    "Riya",
		"lisamccormack85@gmail.com":   "Lisa",
	}
)

func main() {

	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		log.Fatal("GITHUB_TOKEN env not set")
	}

	server := ci_league.NewServer(
		template.Must(template.ParseFiles("template.html")),
		repo,
		token,
		githubEmailMappings,
	)

	http.ListenAndServe(":8000", server)
}
