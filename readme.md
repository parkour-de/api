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

## Prerequisites

You will need Go to run this project. You will also need python3 and pip for the image converter:

    python3 -m pip install -r requirements.txt

## Deployment
This project is set up for Google Cloud Build, as it can build,test the Golang application and push the
resulting Docker image to Google Container Registry.

You can type the `make` command to find all possible targets. Some examples are shown below:

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
By setting the environment variable `PORT`, you can start the server on a different port.
If you provide `UNIX` as an environment variable, it will instead listen on a Unix socket.
See examples below:
```sh
PORT=8080 make run
UNIX=/tmp/dpv.sock make run
```

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
we use `golang:1.22-alpine` as the builder image and `alpine:latest` as the final Docker image. The builder
image includes the necessary tools to build and test the application, while the final image is based on Alpine
Linux, which is a lightweight Linux distribution that is optimized for containerized environments.

The `Dockerfile` contains everything needed to run tests within Docker and build the image. However, when we use
Cloud Build, the testing and building of the binary happens within `cloudbuild.yaml`. Therefore, the
`Dockerfile-cloudbuild` is only responsible for creating the image by copying the binary into it. By separating
the build and image creation steps in this way, we can optimize the build process and make it visible in the UI.

## Setting up ArangoDB

```sh
docker run -p 8529:8529 -e ARANGO_ROOT_PASSWORD=change-me -v /Users/changeme/path/to/dpv-db:/var/lib/arangodb3 arangodb/arangodb:latest
```

Alternatively, you can connect to a server and create a local jump host to connect to the database server. If your db
password is set up accordingly, running the tests will create your ephemeral test databases on the remote server.

```sh
ssh -N -L 8529:127.0.0.1:8529 37.114.34.98
```

## API documentation

**Validate RAML files and generate HTML documentation and JSON file:**

```sh
make raml
```

**Interactive RAML documentation:**

Clone and run [API Console](https://github.com/mulesoft/api-console) from Mulesoft as written in the documentation
and load the demo website for the standalone mode. Choose "Upload an API" in the hamburger menu, choose a .zip file
created from `/docs/api.raml` (make sure to also include the types folder in it) and enjoy!