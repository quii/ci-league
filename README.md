# ci-league ![build status](https://github.com/quii/ci-league/workflows/Test/badge.svg)

Generates a league table of committers to master for a given github repo

## requirements

- Go 1.13

## run

`$ go run cmd/ci-league.go`

Visit `http://localhost:8000?owner={owner}&repo={repo}`

### options via environment variables

- `PORT` defines the port the server listens on (default 8000)
- `GITHUB_TOKEN` to get stats for private repos you'll need a [github access token](https://github.com/settings/tokens)

## test

`$ go test`