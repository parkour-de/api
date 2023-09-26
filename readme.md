# Deutscher Parkour Verband API

The Deutscher Parkour Verband API is a web service that provides access to information related to parkour training
sessions, locations, associations, members, users and more. It's designed to be a part of the infrastructure for
managing and organizing parkour activities within the German Parkour Association.

## Purpose

The primary purpose of this API is to serve as the backend for a comprehensive platform for parkour enthusiasts,
trainers, and organizers. It allows users to:

- Retrieve information about upcoming training sessions.
- Find parkour training locations in their area.
- Register and manage user profiles.
- Organize and schedule training sessions and events.
- Sort, filter, comment and share information.

## Features

This API comes with a variety of features, including:

- **User Management:** Users can register, log in, and manage their profiles.
- **Training Sessions:** Information about upcoming training sessions, including date, time, location, and organizers.
- **Locations:** A database of parkour training locations, complete with details on facilities and accessibility.
- **Event Scheduling:** Organizers can create and schedule training sessions, workshops, and events.
- **OAuth Integration:** Support for OAuth authentication, allowing users to log in with their preferred social media accounts.

## Deployment
This project is set up for Google Cloud Build, as it can build,test the Golang application and push the
resulting Docker image to Google Container Registry.

## Local Testing
To build and test the application locally, run the following command:

```sh
make build
```
This will compile the application and run the unit tests. To run the application locally, use the following command:
```sh
make run
```
This will start the HTTP server at http://localhost:8080.

If you prefer to run the application inside a Docker container, you can use the following command instead:
```sh
make docker-build docker-run
```
This will build a Docker image and start a container that runs the application. The HTTP server will be available
at http://localhost:8080.

## Google Cloud Build
To build and push the Docker image to Google Container Registry using Cloud Build, use the following command:
```sh
gcloud builds submit --config=cloudbuild.yaml .
```
The `cloudbuild.yaml` file contains the build steps for the Cloud Build process. The steps include downloading
Go modules, running unit tests, building the application, building the Docker image, and pushing the image to
Google Container Registry.

In order to optimize the build process and make the resulting Docker image as small and performant as possible,
we use `golang:1.21-alpine` as the builder image and `alpine:latest` as the final Docker image. The builder
image includes the necessary tools to build and test the application, while the final image is based on Alpine
Linux, which is a lightweight Linux distribution that is optimized for containerized environments.

The `Dockerfile` contains everything needed to run tests within Docker and build the image. However, when we use
Cloud Build, the testing and building of the binary happens within `cloudbuild.yaml`. Therefore, the
`Dockerfile-cloudbuild` is only responsible for creating the image by copying the binary into it. By separating
the build and image creation steps in this way, we can optimize the build process and make it visible in the UI.

## Setting up ArangoDB

    docker run -p 8529:8529 -e ARANGO_ROOT_PASSWORD=change-me arangodb/arangodb:3.11.3

## Update Swagger documentation

    go install github.com/swaggo/swag/cmd/swag@latest
    swag fmt -g src/cmd/endpoint1/main.go
    swag init -g src/cmd/endpoint1/main.go
    npm install @redocly/cli -g
    redocly build-docs docs/swagger.yaml  

Then, open the file in a browser and click Download to convert the result to OpenAPI 3.0
This allows the file to be imported in tools like Insomnia.
Some example arrays in the definition files are not correctly formatted and require manual chanages.