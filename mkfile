#!/usr/bin/env bash

set -xe

d:
  go run cmd/app/main.go

b:
  go build -o bin/app cmd/app/main.go

env:
  export USER_ID="$(id -u)"

up: env
  docker compose up -d --build

down: env
  docker compose down -v

db:
  sqlc generate

qlint:
  sqlc vet

migration:
  migrate create -ext sql -dir db/migrations -seq $1

mocks:
  mockery

t:
  go test ./...

lint:
  golangci-lint run ./...
