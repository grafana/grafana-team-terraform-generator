# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH environment variable at /go.
FROM golang:1.27@sha256:5f1b98fa02fe396ebce2a0227a0e08309e11792915b9dd2e1f71053d7e7e6666 as builder

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