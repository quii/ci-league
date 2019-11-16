# ci-league

Generates a league table of committers to master for a given github repo

## requirements

- Go 1.13

## run

`$ go run cmd/ci-league.go`

Visit `http://localhost:8000?owner={owner}&repo={repo}`

## test

`$ go test`