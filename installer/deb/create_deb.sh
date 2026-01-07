#!/bin/bash

# This script generates a .deb file for the installation of MassaStation on a Debian-based Linux distribution.

set -e

BUILD_DIR=builddeb
TMP_DIR=tmpdeb
PKGVERSION=0.0.0-dev

MASSASTATION_ARCHIVE_NAME=MassaStation.tar.xz
MASSASTATION_BINARY_NAME=massastation

# Print error message to stderr and exit with code 1.
fatal() {
    echo "FATAL: $1"
    exit 1
}

# Install dependencies required to build the MassaStation binary.
install_massastation_build_dependencies() {
    sudo apt-get update || fatal "failed to update apt"
    sudo apt-get install -y --fix-missing libgl1-mesa-dev xorg-dev || fatal "failed to install libgl1-mesa-dev xorg-dev"
    go install fyne.io/tools/cmd/fyne@v1.7.0 || fatal "failed to install fyne.io/tools/cmd/fyne@v1.7.0"
    go install github.com/go-swagger/go-swagger/cmd/swagger@latest || fatal "failed to install github.com/go-swagger/go-swagger/cmd/swagger@latest"
    go install golang.org/x/tools/cmd/stringer@latest || fatal "failed to install golang.org/x/tools/cmd/stringer@latest"
}

# Build MassaStation from source.
build_massastation() {
    install_massastation_build_dependencies

    # Ensure go install binaries are in PATH
    export PATH="$PATH:$(go env GOPATH)/bin"

    go generate ../... || fatal "go generate failed for $MASSASTATION_BINARY_NAME"
    export GOARCH=$ARCH
    export CGO_ENABLED=1
    fyne package -src ../cmd/massastation -icon ../../int/systray/embedded/logo.png -name MassaStation --app-id com.massalabs.massastation || fatal "fyne package failed for $MASSASTATION_BINARY_NAME"
}

# Delete the build directory if it exists.
clean() {
    if [ -d $BUILD_DIR ]; then
        rm -rf $BUILD_DIR || fatal "failed to delete $BUILD_DIR"
    fi

    if [ -d $TMP_DIR ]; then
        rm -rf $TMP_DIR || fatal "failed to delete $TMP_DIR"
    fi
}

# Install dependencies required to build the .deb file.
install_dependencies() {
    sudo apt-get install dpkg-dev -y || fatal "failed to install dpkg-dev"
}

main() {
    clean

    install_dependencies

    test -f $MASSASTATION_ARCHIVE_NAME || build_massastation

    mkdir -p $TMP_DIR || fatal "failed to create $TMP_DIR"
    tar -xf $MASSASTATION_ARCHIVE_NAME -C $TMP_DIR || fatal "failed to extract $MASSASTATION_ARCHIVE_NAME to $TMP_DIR"

    mkdir -p $BUILD_DIR/usr/bin || fatal "failed to create $BUILD_DIR/usr/bin"

    # Check if the binary isn't named massastation. If it isn't, rename it to massastation.
    if [ ! -f $TMP_DIR/usr/local/bin/$MASSASTATION_BINARY_NAME ]; then
        mv $TMP_DIR/usr/local/bin/massastation_* $TMP_DIR/usr/local/bin/$MASSASTATION_BINARY_NAME || fatal "failed to rename binary to $MASSASTATION_BINARY_NAME"
    fi
    cp $TMP_DIR/usr/local/bin/$MASSASTATION_BINARY_NAME $BUILD_DIR/usr/bin || fatal "failed to copy $MASSASTATION_BINARY_NAME to $BUILD_DIR/usr/bin"
    chmod +x $BUILD_DIR/usr/bin/$MASSASTATION_BINARY_NAME || fatal "failed to make $MASSASTATION_BINARY_NAME executable"

    mkdir -p $BUILD_DIR/usr/share/applications || fatal "failed to create $BUILD_DIR/usr/share/applications"
    cp $TMP_DIR/usr/local/share/applications/net.massalabs.massastation.desktop $BUILD_DIR/usr/share/applications || fatal "failed to copy net.massalabs.massastation.desktop to $BUILD_DIR/usr/share/applications"

    mkdir -p $BUILD_DIR/usr/share/pixmaps || fatal "failed to create $BUILD_DIR/usr/share/pixmaps"
    cp $TMP_DIR/usr/local/share/pixmaps/net.massalabs.massastation.png $BUILD_DIR/usr/share/pixmaps || fatal "failed to copy net.massalabs.massastation.png to $BUILD_DIR/usr/share/pixmaps"

    mkdir -p $BUILD_DIR/usr/share/doc/massastation || fatal "failed to create $BUILD_DIR/usr/share/doc/massastation"
    cp common/MassaStation_ToS.txt $BUILD_DIR/usr/share/doc/massastation/terms-and-conditions.txt || fatal "failed to copy MassaStation_ToS.txt to $BUILD_DIR/usr/share/doc/massastation/terms-and-conditions.txt"

    mkdir -p $BUILD_DIR/DEBIAN || fatal "failed to create $BUILD_DIR/DEBIAN"
    cat <<EOF >$BUILD_DIR/DEBIAN/control
Package: massastation
Version: $PKGVERSION
Architecture: amd64
Maintainer: Massa Labs <massa.net>
Homepage: https://station.massa.net
Description: An entrance to the Massa blockchain.
    MassaStation is a secured gateway to the Massa blockchain. This application provides a user-friendly way to access, use and build on the Massa blockchain while keeping you safe from the dangers of the internet.
Depends: iproute2, libnss3-tools, debconf (>= 0.5) | debconf-2.0
Recommends: libwebkit2gtk-4.1-dev
EOF

    cp deb/scripts/* $BUILD_DIR/DEBIAN/ || fatal "failed to copy installer scripts to $BUILD_DIR/DEBIAN"
    cp deb/templates $BUILD_DIR/DEBIAN/templates || fatal "failed to copy templates to $BUILD_DIR/DEBIAN"

    DEB_NAME=massastation_$PKGVERSION\_amd64.deb
    dpkg-deb --build $BUILD_DIR $DEB_NAME || fatal "failed to build $DEB_NAME"
}

# Check if $VERSION is set and set $PKGVERSION to $VERSION.
if [ ! -z "$VERSION" ]; then
    # Remove the `v` prefix from the version.
    PKGVERSION=$VERSION
else # If $VERSION is not set, use the latest git tag followed by `-dev`
    PKGVERSION=$(git describe --tags --abbrev=0 --match 'v[0-9]*' | sed 's/^v//')-dev
fi

main
