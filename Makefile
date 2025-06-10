SHELL := /bin/bash

run:
	go run ./main.go

lint:
	golangci-lint run ./...

test_race:
	go test -race -timeout=60s -count 1 ./...

test:
	go clean -testcache
	go test ./...

.PHONY: run lint test_race test