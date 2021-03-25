# simple makefile @avrebarra
NAME=main
COVERAGE_MIN=50.0

## install: setup project
install:
	go mod tidy

## watch: development with air
watch:
	air -c .air.toml

## build: Build binary applications
build:
	go generate ./...
	go build -o ./dist/${NAME} .

## build-linux: Build binary applications for linux
build-linux:
	go generate ./...
	GOOS=linux go build -o ./dist/${NAME} .

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run with parameter options: "
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
