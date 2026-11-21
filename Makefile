.PHONY: dev build run tidy fmt vet clean seed

DB_PATH ?= ./barbershop.db
PORT ?= 8080

dev: ## Run with hot reload (air)
	air

build: ## Build the local server binary
	go build -o bin/localserver .

run: build ## Build and run (DB at ./barbershop.db)
	BARBERHOP_DB=$(DB_PATH) PORT=$(PORT) ./bin/localserver

seed: ## Create initial admin user (interactive)
	go run cmd/seed/main.go $(DB_PATH)

tidy: ## Sync go.mod and go.sum
	go mod tidy

fmt: ## Format the code
	gofmt -w .

vet: ## Static checks
	go vet ./...

clean: ## Remove build artifacts
	rm -rf bin