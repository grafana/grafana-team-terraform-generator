# Use the official Go image to create a build artifact.
FROM golang:1.26@sha256:3aff6657219a4d9c14e27fb1d8976c49c29fddb70ba835014f477e1c70636647 AS builder

WORKDIR /go/src/github.com/grafana/grafana-team-terraform-generator
COPY . .

ENV GOTOOLCHAIN=local
RUN go build -trimpath -ldflags="-s -w" -o /grafana-tf-gen

# Use a lean, supported Debian runtime image.
FROM debian:bookworm-slim@sha256:7b140f374b289a7c2befc338f42ebe6441b7ea838a042bbd5acbfca6ec875818

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /grafana-tf-gen /grafana-tf-gen

CMD ["/grafana-tf-gen"]
