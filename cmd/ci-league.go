package main

import (
	"fmt"
	"github.com/quii/ci-league"
	"github.com/quii/ci-league/github"
	"html/template"
	"log"
	"net/http"
	"os"
)

const defaultPort = ":8000"

func main() {
	client := github.NewClient(os.Getenv("GITHUB_TOKEN"))

	mappings := map[string]string{
		"tamara.jordan1+coding@hotmail.com":                         "Tamara",
		"27856297+dependabot-preview[bot]@users.noreply.github.com": "Depandabot",
		"qui666@gmail.com":            "Chris",
		"riyaddattani@gmail.com":      "Riya",
		"rick@22px.io":                "Ricky",
		"karol.slomczynski@gmail.com": "Osh",
		"riya_dattani@hotmail.com":    "Riya",
		"lisamccormack85@gmail.com":   "Lisa",
		"reis.ivo@gmail.com":          "Ivo",
		"ckurzeja@scottlogic.com":     "CK",
	}

	service := ci_league.NewLeagueService(github.NewService(client), mappings)

	server := ci_league.NewServer(
		template.Must(template.ParseFiles("template.html")),
		service,
	)

	port := getPort()

	fmt.Println("Listening on port", port)
	if err := http.ListenAndServe(port, server); err != nil {
		log.Fatalf("Couldn't launch server listening on %s, %s", port, err)
	}
}

func getPort() string {
	port := defaultPort
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	return port
}
