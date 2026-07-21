# Extends the prebuilt Microsoft Dev Containers Go base image.
# Features (Terraform, GitHub CLI, and common-utils) are layered via devcontainer.json.
FROM mcr.microsoft.com/devcontainers/go:1.26@sha256:638cb8e8c040d29fbd0bb9a0d174b10aef96683437f839048d091eb337a8ea85

# Extra OS-level tools the generator or its workflow may need beyond the features.
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        git \
        make \
        jq \
    && rm -rf /var/lib/apt/lists/*

# Keep local validation reproducible instead of downloading a different toolchain.
ENV GOTOOLCHAIN=local
