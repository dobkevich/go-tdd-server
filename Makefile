.PHONY: build test run lint clean docker-build help

BINARY_NAME=server
GO_VERSION=1.24

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -buildvcs=false -v -o $(BINARY_NAME) main.go

test: ## Run tests
	go test -buildvcs=false -v -race ./...

run: build ## Build and run the binary
	./$(BINARY_NAME)

lint: ## Run golangci-lint
	GOFLAGS="-buildvcs=false" golangci-lint run ./...

clean: ## Remove binary and logs
	rm -f $(BINARY_NAME)
	rm -f *.log

docker-build: ## Build Docker image
	docker build -t go-tdd-server:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 go-tdd-server:latest
