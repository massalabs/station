#!/bin/bash +x

BINARY_URL="https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin"
SCRIPT="MacOS"

MKCERT_URL_ARM="https://dl.filippo.io/mkcert/latest?for=darwin/arm64"
MKCERT_URL_AMD="https://dl.filippo.io/mkcert/latest?for=darwin/amd64"

THYRA_CONF_DIR=$HOME/.config/thyra

green () { echo -e "\033[01;32m$1:\033[0m $2"; }

fatal () { echo -e "\033[01;31m[$SCRIPT]ERROR:\033[0m $1" >&2; exit 1; }

architecture_version () {
    case $(uname -m) in
        x86_64) echo "amd64" ;;
        arm64)  echo "arm64" ;;
        *)      fatal "$(uname -m) unsupported." ;;
    esac
}

install_thyra () {
    arch="$(architecture_version)" || exit 1
    curl -s -L "${BINARY_URL}_${arch}" -o thyra-server || fatal "binary download failed."

    chmod +x thyra-server || fatal "change to executable failed."
    $(sudo [ -d /usr/local/bin ] || sudo mkdir /usr/local/bin) || fatal "/usr/local/bin creation failed."
    sudo mv thyra-server /usr/local/bin/ || fatal "move to /usr/local/bin/ failed."

    # Create config dir
    mkdir -p $THYRA_CONF_DIR
}

configure_start_dnsmasq () {
    sudo bash -c 'echo "address=/.massa/127.0.0.1" > $(brew --prefix)/etc/dnsmasq.d/massa.conf' || fatal "dnsmasq configuration failed."
    sudo mkdir -p /etc/resolver  || fatal "resolver directory creation failed."
    sudo bash -c 'echo "nameserver 127.0.0.1" > /etc/resolver/massa'  || fatal "resolver configuration failed."

    sudo brew services restart dnsmasq || fatal "dnsmasq service failed to start."
}

set_local_dns () {
    case $(sudo lsof -i :53 | sed -n 2p | sed 's/[[:space:]].*$//') in
        "")         (brew install dnsmasq || fatal "dnsmasq installation failed.") && configure_start_dnsmasq || exit -1 ;;
        dnsmasq)    configure_start_dnsmasq || exit -1 ;;
        *)          fatal "local DNS application unsupported." ;;
    esac
}

generate_certificate () {
    green "INFO" "Installing MKcert, its dependencies and then generating HTTPS certificates:"

    [[ $(find /Applications/ -type d -iname "*Firefox*.app") ]] && (brew install nss || fatal "impossible to install certutil. Thyra will not work on Firefox.")
    
    ARCH=$(uname -m)
    if [[ "$ARCH" == 'aarch64' ]]; then
        MKCERT_URL=$MKCERT_URL_ARM
        else
        MKCERT_URL=$MKCERT_URL_AMD
    fi
    curl -sL $MKCERT_URL -o mkcert
    chmod +x mkcert
    ./mkcert -install
    mkdir -p $THYRA_CONF_DIR/certs
    ./mkcert -cert-file $THYRA_CONF_DIR/certs/cert.pem -key-file $THYRA_CONF_DIR/certs/cert-key.pem my.massa
    rm mkcert
}

echo ""

green "INFO" "This installation script will install the last release of Thyra and will configure your local DNS to resolve .massa if needed."

install_thyra || exit 1

generate_certificate || exit 1

ping -c 1 -t 1 test.massa  || set_local_dns || exit 1

green "SUCCESS" "Thyra is installed and the .massa TLD resolution is configured. You're free to go!!!"

echo ""
