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
    sudo apt-get install libgl1-mesa-dev xorg-dev -y || fatal "failed to install libgl1-mesa-dev xorg-dev"
    go install fyne.io/fyne/v2/cmd/fyne@latest || fatal "failed to install fyne.io/fyne/v2/cmd/fyne@latest"
    go install github.com/go-swagger/go-swagger/cmd/swagger@latest || fatal "failed to install github.com/go-swagger/go-swagger/cmd/swagger@latest"
    go install golang.org/x/tools/cmd/stringer@latest || fatal "failed to install golang.org/x/tools/cmd/stringer@latest"
}

# Build MassaStation from source.
build_massastation() {
    install_massastation_build_dependencies

    go generate ../... || fatal "go generate failed for $MASSASTATION_BINARY_NAME"
    export GOARCH=$ARCH
    export CGO_ENABLED=1
    fyne package -icon logo.png -name MassaStation -appID com.massalabs.massastation -src ../cmd/massastation || fatal "fyne package failed for $MASSASTATION_BINARY_NAME"
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
    cp $TMP_DIR/usr/local/share/applications/MassaStation.desktop $BUILD_DIR/usr/share/applications || fatal "failed to copy MassaStation.desktop to $BUILD_DIR/usr/share/applications"

    mkdir -p $BUILD_DIR/usr/share/pixmaps || fatal "failed to create $BUILD_DIR/usr/share/pixmaps"
    cp $TMP_DIR/usr/local/share/pixmaps/MassaStation.png $BUILD_DIR/usr/share/pixmaps || fatal "failed to copy MassaStation.png to $BUILD_DIR/usr/share/pixmaps"

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
