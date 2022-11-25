#!/bin/bash
set -e

generate-certificate() {
    MKCERT_URL="https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-"
    BIN_LINUX="linux-amd64"
    BIN_MACOS_ARM="darwin-arm64"
    BIN_MACOS_AMD="darwin-amd64"

    if [[ "$1" == 'MacOS' ]]; then
        ARCH=$(uname -m)
        if [[ "$ARCH" == 'aarch64' ]]; then
            URL=$MKCERT_URL$BIN_MACOS_ARM
        else
            URL=$MKCERT_URL$BIN_MACOS_AMD
        fi
    else
        URL=$MKCERT_URL$BIN_LINUX
    fi

    curl -sL $URL -o mkcert
    chmod +x mkcert
    ./mkcert -install
    mkdir -p $THYRA_CONF_DIR/certs
    ./mkcert -cert-file $THYRA_CONF_DIR/certs/cert.pem -key-file $THYRA_CONF_DIR/certs/cert-key.pem my.massa
    rm mkcert
}
