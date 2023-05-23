#!/bin/bash

# This script generates a .pkg file for the installation of MassaStation on Mac OS.

set -e

BUILD_DIR=buildpkg
PKGVERSION=dev
ARCH=$1

SERVER_BINARY_NAME=massastation-server
APP_BINARY_NAME=massastation-app

# Print the usage to stderr and exit with code 1.
display_usage() {
    echo "Usage: $0 <arch>" >&2
    echo "  arch: amd64 or arm64" >&2
    exit 1
}

# Print error message to stderr and exit with code 1.
fatal() {
    echo "FATAL: $1"
    exit 1
}

# Download the latest release of MassaStation app.
download_massastation_app() {
    curl -L https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_darwin-$ARCH -o $APP_BINARY_NAME || fatal "failed to download $APP_BINARY_NAME"
    chmod +x $APP_BINARY_NAME || fatal "failed to chmod $APP_BINARY_NAME"
}

# Download the latest release of MassaStation server.
build_massastation_server() {
    go generate ../... || fatal "go generate failed for $SERVER_BINARY_NAME"
    go build -o $SERVER_BINARY_NAME ../cmd/thyra-server || fatal "failed to build $SERVER_BINARY_NAME"
    chmod +x $SERVER_BINARY_NAME || fatal "failed to chmod $SERVER_BINARY_NAME"
}

# Delete the build directory if it exists.
clean() {
    if [ -d $BUILD_DIR ]; then
        rm -rf $BUILD_DIR || fatal "failed to delete $BUILD_DIR"
    fi
}

# Build the package using pkgbuild.
package() {
    pkgbuild --root $BUILD_DIR --identifier com.massalabs.massastation --version $PKGVERSION \
        --scripts macos/scripts --install-location / massastation_$PKGVERSION\_$ARCH.pkg || fatal "failed to create package"
}

main() {
    clean

    test -f $SERVER_BINARY_NAME || build_massastation_server
    test -f $APP_BINARY_NAME || download_massastation_app

    mkdir -p $BUILD_DIR/usr/local/bin || fatal "failed to create $BUILD_DIR/usr/local/bin"
    cp $SERVER_BINARY_NAME $BUILD_DIR/usr/local/bin || fatal "failed to copy $SERVER_BINARY_NAME to $BUILD_DIR/usr/local/bin"
    cp $APP_BINARY_NAME $BUILD_DIR/usr/local/bin || fatal "failed to copy $APP_BINARY_NAME to $BUILD_DIR/usr/local/bin"

    package
}

# Check if the user provided the build architecture.
if [ -z "$ARCH" ]; then
    display_usage
fi

# Check if the user provided a valid build architecture.
if [ "$ARCH" != "amd64" ] && [ "$ARCH" != "arm64" ]; then
    display_usage
fi

# Check if $VERSION is set and set $PKGVERSION to $VERSION.
if [ ! -z "$VERSION" ]; then
    PKGVERSION=$VERSION
else # If $VERSION is not set, use the latest git tag followed by `-dev`
    PKGVERSION=$(git describe --tags --abbrev=0 | sed 's/^v//')-dev
fi

main
