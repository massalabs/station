#!/bin/bash

# This script generates a .pkg file for the installation of MassaStation on Mac OS.

set -e

PKGVERSION=dev
ARCH=$1

MASSASTATION_INSTALLER_NAME=MassaStation.app
MASSASTATION_BINARY_NAME=massastation

HOMEBREW_INSTALL_SCRIPT_URL=https://raw.githubusercontent.com/massalabs/homebrew.sh/master/homebrew-3.3.sh

LICENSE_FILE_NAME=MassaStation_ToS.rtf

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

# Install dependencies required to build the MassaStation binary.
install_massastation_build_dependencies() {
    go install fyne.io/fyne/v2/cmd/fyne@latest
    go install github.com/go-swagger/go-swagger/cmd/swagger@latest
    go install golang.org/x/tools/cmd/stringer@latest
}

# Build MassaStation from source.
build_massastation() {
    install_massastation_build_dependencies

    go generate ../... || fatal "go generate failed for $MASSASTATION_INSTALLER_NAME"
    export GOARCH=$ARCH
    export CGO_ENABLED=1
    # -icon is based on the path of the -src flag.
    fyne package -icon ../../int/systray/embedded/logo.png -name MassaStation -appID com.massalabs.massastation -src ../cmd/massastation || fatal "fyne package failed for $MASSASTATION_INSTALLER_NAME"
    chmod +x $MASSASTATION_INSTALLER_NAME || fatal "failed to chmod $MASSASTATION_INSTALLER_NAME"
}

# Build the package using pkgbuild.
package() {
    pkgbuild --component $MASSASTATION_INSTALLER_NAME --identifier com.massalabs.massastation --version $PKGVERSION \
        --scripts macos/scripts --install-location /Applications MassaStation.pkg || fatal "failed to create package"
    
    productbuild --distribution macos/Distribution.dist --resources macos/resources --package-path . massastation_$PKGVERSION\_$ARCH.pkg || fatal "failed to create installer"
}

# Download homebrew installation script and put it in script directory.
download_homebrew_install_script() {
    curl -sL $HOMEBREW_INSTALL_SCRIPT_URL -o macos/scripts/install_homebrew.sh || fatal "failed to download homebrew installation script"
    chmod +x macos/scripts/install_homebrew.sh || fatal "failed to chmod homebrew installation script"
}

main() {
    test -d $MASSASTATION_INSTALLER_NAME || build_massastation

    download_homebrew_install_script

    if [ ! -d macos/resources ]; then
        mkdir macos/resources || fatal "failed to create resources directory"
    fi

    if [ ! -f macos/resources/$LICENSE_FILE_NAME ]; then
        cp common/$LICENSE_FILE_NAME macos/resources/$LICENSE_FILE_NAME || fatal "failed to copy license file"
    fi

    # Check if the binary isn't named massastation. If it isn't, rename it to massastation.
    if [ ! -f $MASSASTATION_INSTALLER_NAME/Contents/MacOS/$MASSASTATION_BINARY_NAME ]; then
        mv $MASSASTATION_INSTALLER_NAME/Contents/MacOS/massastation_* $MASSASTATION_INSTALLER_NAME/Contents/MacOS/$MASSASTATION_BINARY_NAME || fatal "failed to rename $MASSASTATION_INSTALLER_NAME to $MASSASTATION_BINARY_NAME"
    fi

    chmod +x $MASSASTATION_INSTALLER_NAME/Contents/MacOS/$MASSASTATION_BINARY_NAME || fatal "failed to chmod $MASSASTATION_INSTALLER_NAME/Contents/MacOS/$MASSASTATION_BINARY_NAME"

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
    PKGVERSION=$(echo $VERSION | sed 's/^v//')
else # If $VERSION is not set, use the latest git tag followed by `-dev`
    PKGVERSION=$(git describe --tags --abbrev=0 | sed 's/^v//')-dev
fi

main
