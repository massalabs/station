#!/bin/bash

# This script generates a .deb file for the installation of MassaStation on a Debian-based Linux distribution.

BUILD_DIR=builddeb
PKGVERSION=dev

DEB_NAME=massastation_$PKGVERSION\_amd64.deb
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

# Download the latest release of MassaStation server.
download_massastation_server() {
    curl -L https://github.com/massalabs/thyra/releases/latest/download/thyra-server_linux_amd64 -o $SERVER_BINARY_NAME || fatal "failed to download $SERVER_BINARY_NAME"
    chmod +x $SERVER_BINARY_NAME || fatal "failed to chmod $SERVER_BINARY_NAME"
}

# Delete the build directory if it exists.
clean() {
    test -d $BUILD_DIR && rm -rf $BUILD_DIR
}

# Install dependencies required to build the .deb file.
install_dependencies() {
    sudo apt-get install dpkg-dev
}

main() {
    clean

    install_dependencies

    test -f $SERVER_BINARY_NAME || download_massastation_server
    test -f $APP_BINARY_NAME || download_massastation_app

    mkdir -p $BUILD_DIR/usr/bin
    cp $SERVER_BINARY_NAME $BUILD_DIR/usr/bin
    cp $APP_BINARY_NAME $BUILD_DIR/usr/bin

    mkdir -p $BUILD_DIR/DEBIAN
    cat <<EOF >$BUILD_DIR/DEBIAN/control
Package: massastation
Version: $PKGVERSION
Architecture: amd64
Maintainer: Massa Labs <massa.net>
Homepage: https://github.com/massalabs/thyra
Description: An entrance to the Massa blockchain.
    MassaStation is a secured gateway to the Massa blockchain. 
    This application provides a user-friendly way to access, use and build on the Massa blockchain while
    keeping you safe from the dangers of the internet.
Recommends: libnss3-tools
EOF

    cp deb/scripts/postinst $BUILD_DIR/DEBIAN

    dpkg-deb --build $BUILD_DIR massastation_$PKGVERSION\_amd64.deb
}

# Check if $VERSION is set and set $PKGVERSION to $VERSION.
if [ ! -z "$VERSION" ]; then
    PKGVERSION=$VERSION
fi

main
