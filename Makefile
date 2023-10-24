# Variables
IMAGE_NAME=dpv
CONTAINER_NAME=dpv_container
PORT=8080
OUT_DIR := ./bin

# Targets
.DEFAULT_GOAL:=help
.PHONY: all build test run docker-build docker-run docker-stop help

all: test build ## Run test, then build

build: ## Build the binary
	go build -o $(OUT_DIR)/endpoint1 -ldflags "-X main.version=$$(git rev-list --count main)" ./src/cmd/endpoint1

test: ## Run tests
	go test ./... -p 8

run: ## Run the binary
	$(OUT_DIR)/endpoint1

docker-build: ## Build the container image
	docker build -t $(IMAGE_NAME) --build-arg VERSION=$$(git rev-list --count main) .

docker-run: ## Run the container
	docker run -d -p $(PORT):$(PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME)

docker-stop: ## Stop the container
	docker stop $(CONTAINER_NAME) && docker rm $(CONTAINER_NAME)

raml: ## Update raml documentation
	(npm list -g raml2html && npm list -g oas-raml-converter-cli) || npm i -g raml2html oas-raml-converter-cli
	raml2html -v -i docs/api.raml -o docs/api.html
	expect ./oas-raml-converter-cli.expect

help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
