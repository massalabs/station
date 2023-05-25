#!/bin/bash

# This script generates a .pkg file for the installation of MassaStation on Mac OS.

set -e

PKGVERSION=dev
ARCH=$1

MASSASTATION_BINARY_NAME=MassaStation.app

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

# Build MassaStation from source.
build_massastation() {
    go generate ../... || fatal "go generate failed for $MASSASTATION_BINARY_NAME"
    export GOARCH=$ARCH
    export CGO_ENABLED=1
    fyne package -icon logo.png -name MassaStation -appID com.massalabs.massastation -src ../cmd/massastation || fatal "fyne package failed for $MASSASTATION_BINARY_NAME"
    chmod +x $MASSASTATION_BINARY_NAME || fatal "failed to chmod $MASSASTATION_BINARY_NAME"
}

# Build the package using pkgbuild.
package() {
    pkgbuild --component $MASSASTATION_BINARY_NAME --identifier com.massalabs.massastation --version $PKGVERSION \
        --scripts macos/scripts --install-location /Applications massastation_$PKGVERSION\_$ARCH.pkg || fatal "failed to create package"
}

download_dependencies() {
    go get fyne.io/fyne/v2@latest
    go install fyne.io/fyne/v2/cmd/fyne@latest
}

main() {
    download_dependencies

    test -d $MASSASTATION_BINARY_NAME || build_massastation

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
