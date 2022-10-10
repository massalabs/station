#!/bin/bash +x

BINARY_URL="https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin"
SCRIPT="MacOS"

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
    sudo mv thyra-server /usr/local/bin/ || fatal "move to /usr/local/bin/ failed."
}

configure_start_dnsmasq () {
    sudo bash -c 'echo "address=/.massa/127.0.0.1" > $(brew --prefix)/etc/dnsmasq.d/massa.conf' || fatal "dnsmas configuration failed."
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

echo ""

green "INFO" "This installation script will install the last release of Thyra and will configure your local DNS to resolve .massa if needed."

install_thyra || exit 1

ping -c test.massa || set_local_dns || exit 1

green "SUCCESS" "Thyra is installed and the .massa TLD resolution is configured. You're free to go!!!"

echo ""
