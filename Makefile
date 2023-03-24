# Variables
IMAGE_NAME=myapp
CONTAINER_NAME=myapp_container
PORT=8080

# Targets
.PHONY: all
all: test build

build:
	go build -o server .

test:
	go test ./...

run:
	./server

docker-build:
	DOCKER_CONFIG=~/invalid/ docker build -t $(IMAGE_NAME) .

docker-run:
	docker run -d -p $(PORT):$(PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME)

docker-stop:
	docker stop $(CONTAINER_NAME) && docker rm $(CONTAINER_NAME)