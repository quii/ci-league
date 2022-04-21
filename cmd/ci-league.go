package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/quii/ci-league/github"
	"github.com/quii/ci-league/league"
)

const defaultPort = ":8000"
const templatePath = "template.html"

func main() {
	ghToken, tokenSet := os.LookupEnv("GITHUB_TOKEN")
	if !tokenSet {
		fmt.Println("GITHUB_TOKEN not set")
	}

	client := github.NewClient(ghToken, os.Stderr)

	service := github.NewService(client)
	//service := github.NewCachedService(newService, os.Stdout)

	server := league.NewServer(
		template.Must(template.ParseFiles(templatePath)),
		league.NewService(service, getMappings()),
	)

	port := getPort()

	fmt.Println("Listening on port", port)
	fmt.Printf("Try http://localhost%s/integrations?owner=quii&repo=ci-league\n", port)
	if err := http.ListenAndServe(port, server); err != nil {
		log.Fatalf("Couldn't launch server listening on %s, %s", port, err)
	}
}

func getMappings() league.AliasService {
	var mappings league.AliasService
	if os.Getenv("MAPPINGS") != "" {
		m, err := NewFileSystemAliasService(os.Getenv("MAPPINGS"))
		if err != nil {
			log.Fatalf(err.Error())
		}
		mappings = m
	} else {
		mappings = &NoOpAliasService{}
	}
	return mappings
}

func getPort() string {
	port := defaultPort
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}
	return port
}
