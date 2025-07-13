#!/bin/bash

APP_NAME="query_json"
VERSION="1.0.0"

# Get build information
BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Create builds directory
mkdir -p builds

echo "Building ${APP_NAME} v${VERSION} for multiple platforms..."
echo "  Version: ${VERSION}"
echo "  Commit: ${GIT_COMMIT}"
echo "  Build Date: ${BUILD_DATE}"

# Build flags with version information
LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${GIT_COMMIT} -X main.date=${BUILD_DATE}"

# Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o builds/${APP_NAME}-windows-amd64.exe
GOOS=windows GOARCH=386 go build -ldflags="${LDFLAGS}" -o builds/${APP_NAME}-windows-386.exe

# macOS
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o builds/${APP_NAME}-macos-amd64
GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o builds/${APP_NAME}-macos-arm64

# Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o builds/${APP_NAME}-linux-amd64
GOOS=linux GOARCH=386 go build -ldflags="${LDFLAGS}" -o builds/${APP_NAME}-linux-386

echo "Build completed! Binaries are in ./builds/"
ls -la builds/

# Generate checksums
echo "Generating checksums..."
cd builds
sha256sum * > checksums.sha256
echo "Checksums generated in builds/checksums.sha256"
cd ..