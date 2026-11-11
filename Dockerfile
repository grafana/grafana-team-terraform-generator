# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH environment variable at /go.
FROM golang:1.25@sha256:699337d620559a59b4a2bb298ad59611e535d2ee755a34cf2d2a98f37578dc80 as builder

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