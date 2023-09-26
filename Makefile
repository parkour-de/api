# Variables
IMAGE_NAME=dpv
CONTAINER_NAME=dpv_container
PORT=8080

# Targets
.PHONY: all
all: test build

OUT_DIR := ./bin

build:
	go build -o $(OUT_DIR)/endpoint1 ./src/cmd/endpoint1

test:
	go test ./...

run:
	./bin/endpoint1

docker-build:
	DOCKER_CONFIG=~/invalid/ docker build -t $(IMAGE_NAME) .

docker-run:
	docker run -d -p $(PORT):$(PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME)

docker-stop:
	docker stop $(CONTAINER_NAME) && docker rm $(CONTAINER_NAME)