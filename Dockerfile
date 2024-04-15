# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.21 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
# -o myapp specifies the output name of the binary.
RUN GOOS=linux go build -v -ldflags='-s -w -extldflags "-static"' -o myapp

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/
FROM alpine:latest  

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/myapp /myapp

# Run the web service on container startup.
CMD ["/myapp"]
