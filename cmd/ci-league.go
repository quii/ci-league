package main

import (
	"fmt"
	"github.com/quii/ci-league/github"
	"github.com/quii/ci-league/league"
	"html/template"
	"log"
	"net/http"
	"os"
)

const defaultPort = ":8000"
const templatePath = "template.html"

func main() {
	client := github.NewClient(os.Getenv("GITHUB_TOKEN"), os.Stderr)

	mappings := InMemoryAliasService{
		"tamara.jordan1+coding@hotmail.com":                         "Tamara",
		"27856297+dependabot-preview[bot]@users.noreply.github.com": "Depandabot",
		"qui666@gmail.com":                                    "Chris",
		"quii666@gmail.com":                                   "Chris",
		"riyaddattani@gmail.com":                              "Riya",
		"rick@22px.io":                                        "Ricky",
		"karol.slomczynski@gmail.com":                         "Osh",
		"43116906+riyadattani@users.noreply.github.com":       "Riya",
		"lisamccormack85@gmail.com":                           "Lisa",
		"reis.ivo@gmail.com":                                  "Ivo",
		"1410207+ivoreis@users.noreply.github.com":            "Ivo",
		"ckurzeja@scottlogic.com":                             "CK",
		"jmwright@scottlogic.com":                             "Jonny",
		"44469595+jonnymwright@users.noreply.github.com 	": "Jonny",
		"francesco.lentini@gmail.com":                         "Fran",
	}

	service := github.NewService(client)
	//service := github.NewCachedService(newService, os.Stdout)

	server := league.NewServer(
		template.Must(template.ParseFiles(templatePath)),
		league.NewService(service, mappings),
	)

	port := getPort()

	fmt.Println("Listening on port", port)
	fmt.Printf("Try http://localhost%s/integrations?owner=quii&repo=ci-league\n", port)
	if err := http.ListenAndServe(port, server); err != nil {
		log.Fatalf("Couldn't launch server listening on %s, %s", port, err)
	}
}

func getPort() string {
	port := defaultPort
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}
	return port
}
