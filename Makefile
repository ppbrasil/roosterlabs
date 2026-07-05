.PHONY: run test lint build docker

run: ## roda o servidor local em :8080
	go run ./cmd/server

test: ## build + vet + testes (o mesmo que o CI roda)
	go build ./...
	go vet ./...
	go test ./...

lint: ## golangci-lint (instalado no devcontainer/CI)
	golangci-lint run

build:
	CGO_ENABLED=0 go build -trimpath -o bin/server ./cmd/server

docker:
	docker build -t roosterlabs-server .
