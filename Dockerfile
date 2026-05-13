# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH environment variable at /go.
FROM golang:1.24@sha256:d2d2bc1c84f7e60d7d2438a3836ae7d0c847f4888464e7ec9ba3a1339a1ee804 as builder

# Copy the local package files to the container's workspace.
WORKDIR /go/src/github.com/grafana/grafana-terraform-generator
COPY . .

# Build the command inside the container.
# (You might need to modify the path or add additional build commands depending on your app)
RUN go build -o /grafana-tf-gen

# Use a Docker multi-stage build to create a lean production image.
FROM debian:buster-slim

# Copy the binary to the production image from the builder stage.
COPY --from=builder /grafana-tf-gen /grafana-tf-gen

# Run the web service on container startup.
CMD ["/grafana-tf-gen"]