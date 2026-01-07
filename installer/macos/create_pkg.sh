#!/bin/bash

# This script generates a .pkg file for the installation of MassaStation on Mac OS.

# This script can be used in two contexts:
# - in the CI,
# - in local,
# that's why there are some if/else statements.

set -e

PKGVERSION=dev
ARCH=$1

MASSASTATION_APPLICATION_NAME=MassaStation.app
MASSASTATION_BINARY_NAME=massastation

APPLE_DEVELOPER_ID_APPLICATION=$2
APPLE_DEVELOPER_ID_INSTALLER=$3

LICENSE_FILE_NAME=MassaStation_ToS.rtf

# Print the usage to stderr and exit with code 1.
display_usage() {
    echo "Usage: $0 <arch> <APPLE_DEVELOPER_ID_APPLICATION> <APPLE_DEVELOPER_ID_INSTALLER>" >&2
    echo "  arch: amd64 or arm64" >&2
    echo "  APPLE_DEVELOPER_ID_APPLICATION: optional, to sign the .app" >&2
    echo "  APPLE_DEVELOPER_ID_INSTALLER: optional, to sign the .pkg" >&2
    exit 1
}

# Print error message to stderr and exit with code 1.
fatal() {
    echo "FATAL: $1"
    exit 1
}

# Install dependencies required to build the MassaStation binary.
install_massastation_build_dependencies() {
    go install fyne.io/tools/cmd/fyne@v1.7.0
    go install github.com/go-swagger/go-swagger/cmd/swagger@latest
    go install golang.org/x/tools/cmd/stringer@latest
}

# Build MassaStation from source.
build_massastation() {
    install_massastation_build_dependencies

    go generate ../... || fatal "go generate failed for $MASSASTATION_APPLICATION_NAME"
    export GOARCH=$ARCH
    export CGO_ENABLED=1
    fyne package -src ../cmd/massastation -icon ../../int/systray/embedded/logo.png -name MassaStation --app-id com.massalabs.massastation || fatal "fyne package failed for $MASSASTATION_APPLICATION_NAME"
    chmod +x $MASSASTATION_APPLICATION_NAME || fatal "failed to chmod $MASSASTATION_APPLICATION_NAME"
}

# Build the package using pkgbuild.
package() {
    # sign the application if we have a developer id
    if [[ -n "$APPLE_DEVELOPER_ID_APPLICATION" ]]; then
        codesign --force --options runtime --sign "$APPLE_DEVELOPER_ID_APPLICATION" $MASSASTATION_APPLICATION_NAME || fatal "failed to sign"
    fi

    pkgbuild --component $MASSASTATION_APPLICATION_NAME --identifier com.massalabs.massastation --version $PKGVERSION \
        --scripts macos/scripts --install-location /Applications MassaStation.pkg || fatal "failed to create package"

    if [[ -n "$APPLE_DEVELOPER_ID_INSTALLER" ]]; then
        productbuild --distribution macos/Distribution.dist --resources macos/resources --package-path . \
            --sign "$APPLE_DEVELOPER_ID_INSTALLER" massastation_$PKGVERSION\_$ARCH.pkg || fatal "failed to create installer"
    else
        productbuild --distribution macos/Distribution.dist --resources macos/resources --package-path . \
            massastation_$PKGVERSION\_$ARCH.pkg || fatal "failed to create installer"
    fi
}

main() {
    # build massastation only if the .app is not present
    test -d $MASSASTATION_APPLICATION_NAME || build_massastation

    if [ ! -d macos/resources ]; then
        mkdir macos/resources || fatal "failed to create resources directory"
    fi

    if [ ! -f macos/resources/$LICENSE_FILE_NAME ]; then
        cp common/$LICENSE_FILE_NAME macos/resources/$LICENSE_FILE_NAME || fatal "failed to copy license file"
    fi

    # Check if the binary isn't named massastation. If it isn't, rename it to massastation.
    if [ ! -f $MASSASTATION_APPLICATION_NAME/Contents/MacOS/$MASSASTATION_BINARY_NAME ]; then
        mv $MASSASTATION_APPLICATION_NAME/Contents/MacOS/massastation_* $MASSASTATION_APPLICATION_NAME/Contents/MacOS/$MASSASTATION_BINARY_NAME || fatal "failed to rename $MASSASTATION_APPLICATION_NAME to $MASSASTATION_BINARY_NAME"
    fi

    chmod +x $MASSASTATION_APPLICATION_NAME/Contents/MacOS/$MASSASTATION_BINARY_NAME || fatal "failed to chmod $MASSASTATION_APPLICATION_NAME/Contents/MacOS/$MASSASTATION_BINARY_NAME"

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
    # Remove the `v` prefix from the version.
    PKGVERSION=$VERSION
else # If $VERSION is not set, use the latest git tag followed by `-dev`
    PKGVERSION=$(git describe --tags --abbrev=0 --match 'v[0-9]*' | sed 's/^v//')-dev
fi

main
