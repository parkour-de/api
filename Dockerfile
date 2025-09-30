# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.25-alpine as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Run the tests.
RUN go test ./... -p 8

# Build the binary.
ARG VERSION
RUN go build -v -o /app/bin/endpoint1 -ldflags "-X main.version=${VERSION}" ./src/cmd/endpoint1

FROM alpine:latest

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/bin/endpoint1 /app/bin/endpoint1

# Run the web service on container startup.
CMD ["/app/bin/endpoint1"]