#!/bin/bash

set -e

# Get rid of existing binaries
rm -rf dist/ssm*

# Check if VERSION variable set and not empty, otherwise set to default value
if [ -z "$VERSION" ]; then
  VERSION="0.0.1-dev"
fi
echo "Building application version $VERSION"

# Build binaries for different platforms/architectures
OS=linux
ARCH=amd64
echo "Building binaries for $OS/$ARCH..."
GOARCH=$ARCH GOOS=$OS CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/ssm.exe" ssm.go
echo "Building binaries for linux/$ARCH..."
GOARCH=$ARCH GOOS=linux CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/ssm" ssm.go
