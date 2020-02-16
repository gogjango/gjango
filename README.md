![CircleCI](https://img.shields.io/circleci/build/github/calvinchengx/gin-go-pg/master) [![Maintainability](https://api.codeclimate.com/v1/badges/62185b640652168fe9f9/maintainability)](https://codeclimate.com/github/calvinchengx/gin-go-pg/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/62185b640652168fe9f9/test_coverage)](https://codeclimate.com/github/calvinchengx/gin-go-pg/test_coverage) [![Go Report Card](https://goreportcard.com/badge/github.com/calvinchengx/gin-go-pg)](https://goreportcard.com/report/github.com/calvinchengx/gin-go-pg) ![GitHub](https://img.shields.io/github/license/calvinchengx/gin-go-pg)


# golang gin with go-pg orm

An example project that uses golang gin as webserver, and go-pg library for connecting with a PostgreSQL database.

## Get started

```bash
# postgresql config
cp .env.sample .env
source .env
```

```bash
# get dependencies and run
go get -v ./...
go run .
```

## Tests and coverage

```bash
go test -coverprofile c.out ./...
go tool cover -html=c.out
```

## Schema migration

```bash
# either
go run migration/*.go [command]

# or
go build -o bin/migration migration/main.go
./bin/migration [command]

# see all commands via `go run migration/*.go -h`
# we can create migrations by creating files named 1_name.go etc
# reference https://github.com/go-pg/migrations/tree/master/example
```