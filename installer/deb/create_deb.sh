#!/bin/bash

# This script generates a .deb file for the installation of MassaStation on a Debian-based Linux distribution.

set -e

BUILD_DIR=builddeb
PKGVERSION=0.0.0-dev

MASSASTATION_BINARY_NAME=massastation

# Print error message to stderr and exit with code 1.
fatal() {
    echo "FATAL: $1"
    exit 1
}

# Build MassaStation from source.
build_massastation_server() {
    go generate ../... || fatal "go generate failed for $MASSASTATION_BINARY_NAME"
    export GOARCH=$ARCH
    export CGO_ENABLED=1
    fyne package -icon logo.png -name MassaStation -appID com.massalabs.massastation -src ../cmd/massastation || fatal "fyne package failed for $MASSASTATION_BINARY_NAME"
    chmod +x $MASSASTATION_BINARY_NAME || fatal "failed to chmod $MASSASTATION_BINARY_NAME"
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

    test -f $MASSASTATION_BINARY_NAME || build_massastation_server

    mkdir -p $BUILD_DIR/usr/bin || fatal "failed to create $BUILD_DIR/usr/bin"
    cp $MASSASTATION_BINARY_NAME $BUILD_DIR/usr/bin || fatal "failed to copy $MASSASTATION_BINARY_NAME to $BUILD_DIR/usr/bin"

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
