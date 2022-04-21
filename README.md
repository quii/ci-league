# ci-league ![build status](https://github.com/quii/ci-league/workflows/Test/badge.svg)

Generates a league table of committers to master for a given github repo.

[For a longer rant, see my post on CI](https://quii.dev/Gamifying_Continuous_Integration)

## requirements

- Go 1.13

## run

`$ go run cmd/*.go`

or

`$ docker run -p 8000:8000 quii/ci-league`

Visit `http://localhost:8000?owner={owner}&repo={repo}`

It doesn't do any fancy auto-refreshing but most browsers have extensions to auto refresh a tab, we set it for every 10 minutes here.

### options via environment variables

- `PORT` defines the port the server listens on (default 8000)
- `GITHUB_TOKEN` to get stats for private repos you'll need a [github access token](https://github.com/settings/tokens)
- `MAPPINGS` a path to a JSON file with keys of email addresses to aliases
## test

`$ go test ./...`
