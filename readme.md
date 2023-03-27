# Cloud Build Demo Project
This is a demo project for Google Cloud Build that shows how to build and test a Golang application, and push the
resulting Docker image to Google Container Registry.

The application consists of a simple HTTP server that returns "Hello, World!" at the root URL.

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
we use `golang:1.20-alpine` as the builder image and `alpine:latest` as the final Docker image. The builder
image includes the necessary tools to build and test the application, while the final image is based on Alpine
Linux, which is a lightweight Linux distribution that is optimized for containerized environments.

The `Dockerfile` contains everything needed to run tests within Docker and build the image. However, when we use
Cloud Build, the testing and building of the binary happens within `cloudbuild.yaml`. Therefore, the
`Dockerfile-cloudbuild` is only responsible for creating the image by copying the binary into it. By separating
the build and image creation steps in this way, we can optimize the build process and make it visible in the UI.