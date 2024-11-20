HASH ?= $(shell git rev-parse --short HEAD)
VERSION ?= $(shell git describe --tags 2>/dev/null)
PROJECT_NAME ?= $(shell basename "$(PWD)")

all:: help

help ::
	@grep -E '^[a-zA-Z_-]+\s*:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

install :: ## Initial service setup
	@echo "  > Configuring environment from dist files"
	@cp .env.dist .env
	@cp docker-compose.override.yml.dist docker-compose.override.yml
	@echo "  > Installing dependencies"
	@go mod vendor

run :: ## Run service
	docker-compose up -d
	@echo "  > Starting service"
	@go run cmd/server/main.go


format :: ## Formats and standarizes code following conventions
	@echo "  > Formating"
	@gofmt -w .

test :: ## Execute test suite
	@echo "  > Executing tests"
	@go mod vendor

	@echo "  > Running unit"
	@go test -count=1 -cover -race ./...
