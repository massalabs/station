#!/bin/bash

# This script generates a .deb file for the installation of MassaStation on a Debian-based Linux distribution.

set -e

BUILD_DIR=builddeb
PKGVERSION=0.0.0-dev

SERVER_BINARY_NAME=massastation-server
APP_BINARY_NAME=massastation-app

# Print error message to stderr and exit with code 1.
fatal() {
    echo "FATAL: $1"
    exit 1
}

# Download the latest release of MassaStation app.
download_massastation_app() {
    curl -L https://github.com/massalabs/Thyra-Menu-Bar-App/releases/latest/download/ThyraApp_linux-amd64 -o $APP_BINARY_NAME || fatal "failed to download $APP_BINARY_NAME"
    chmod +x $APP_BINARY_NAME || fatal "failed to chmod $APP_BINARY_NAME"
}

# Build the MassaStation server binary.
build_massastation_server() {
    go generate ../... || fatal "go generate failed for $SERVER_BINARY_NAME"
    go build -o $SERVER_BINARY_NAME ../cmd/thyra-server/ || fatal "failed to build $SERVER_BINARY_NAME"
    chmod +x $SERVER_BINARY_NAME || fatal "failed to chmod $SERVER_BINARY_NAME"
}

# Delete the build directory if it exists.
clean() {
    if [ -d $BUILD_DIR ]; then
        rm -rf $BUILD_DIR || fatal "failed to delete $BUILD_DIR"
    fi
}

# Install dependencies required to build the .deb file.
install_dependencies() {
    sudo apt-get install dpkg-dev || fatal "failed to install dpkg-dev"
}

main() {
    clean

    install_dependencies

    test -f $SERVER_BINARY_NAME || build_massastation_server
    test -f $APP_BINARY_NAME || download_massastation_app

    mkdir -p $BUILD_DIR/usr/bin || fatal "failed to create $BUILD_DIR/usr/bin"
    cp $SERVER_BINARY_NAME $BUILD_DIR/usr/bin || fatal "failed to copy $SERVER_BINARY_NAME to $BUILD_DIR/usr/bin"
    cp $APP_BINARY_NAME $BUILD_DIR/usr/bin || fatal "failed to copy $APP_BINARY_NAME to $BUILD_DIR/usr/bin"

    mkdir -p $BUILD_DIR/DEBIAN || fatal "failed to create $BUILD_DIR/DEBIAN"
    cat <<EOF >$BUILD_DIR/DEBIAN/control
Package: massastation
Version: $PKGVERSION
Architecture: amd64
Maintainer: Massa Labs <massa.net>
Homepage: https://github.com/massalabs/thyra
Description: An entrance to the Massa blockchain.
    MassaStation is a secured gateway to the Massa blockchain. This application provides a user-friendly way to access, use and build on the Massa blockchain while keeping you safe from the dangers of the internet.
Recommends: libnss3-tools
EOF

    cp deb/scripts/postinst $BUILD_DIR/DEBIAN
    DEB_NAME=massastation_$PKGVERSION\_amd64.deb

    dpkg-deb --build $BUILD_DIR $DEB_NAME || fatal "failed to build $DEB_NAME"
}

# Check if $VERSION is set and set $PKGVERSION to $VERSION.
if [ ! -z "$VERSION" ]; then
    PKGVERSION=$VERSION
else # If $VERSION is not set, use the latest git tag followed by `-dev`
    PKGVERSION=$(git describe --tags --abbrev=0 | sed 's/^v//')-dev
fi

main
