FROM mcr.microsoft.com/devcontainers/go:1-1.24-bookworm

# Base image ships a stale apt source for the deprecated dl.yarnpkg.com repo
# whose signing key is no longer available, breaking apt-get update for any
# feature (e.g. docker-in-docker) installed on top of this image.
RUN rm -f /etc/apt/sources.list.d/yarn.list
