# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH environment variable at /go.
FROM golang:1.27@sha256:3680233e3204827fbdc66088528ae6d4b3d034f51d03a99d454f6de034888244 as builder

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