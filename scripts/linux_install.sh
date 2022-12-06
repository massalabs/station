#!/bin/bash +x

BINARY_URL="https://github.com/massalabs/thyra/releases/latest/download/thyra-server_linux"
SCRIPT="Linux"

MKCERT_URL="https://dl.filippo.io/mkcert/latest?for=linux/amd64"

THYRA_CONF_DIR=$HOME/.config/thyra

green () { echo -e "\033[01;32m$1:\033[0m $2"; }

warn () { echo -e "\033[01;33m[$SCRIPT]WARNING:\033[0m $1"; }

fatal () { echo -e "\033[01;31m[$SCRIPT]ERROR:\033[0m $1" >&2; exit 1; }

architecture_version () {
    case $(uname -m) in
        x86_64) echo "amd64" ;;
        # arm64)  echo "arm64" ;; ARM64 is unsupported on Linux for now
        *)      fatal "$(uname -m) unsupported." ;;
    esac
}

check_supported_distro () {
    grep '^ID.*=.*ubuntu' /etc/os-release &> /dev/null || fatal "Aborting: You are using an unsupported distribution."
}

install_thyra () {
    arch="$(architecture_version)" || exit 1
    curl -s -L "${BINARY_URL}_${arch}" -o thyra-server || fatal "binary download failed."

    chmod +x thyra-server || fatal "change to executable failed."
    sudo mv thyra-server /usr/local/bin/ || fatal "move to /usr/local/bin/ failed."

    # Create config dir
    mkdir -p $THYRA_CONF_DIR
}

configure_network_manager () {
    sudo cp /etc/NetworkManager/NetworkManager.conf /etc/NetworkManager/NetworkManager.conf_backup_thyra_install || fatal "couldn't create NetworkManager backup"
    local dns="$(grep '^dns=' /etc/NetworkManager/NetworkManager.conf | sed 's/^dns=//')"
    case $dns in
        dnsmasq);;
        "") sudo sed -i 's/^\[main\]$/\[main\]\ndns=dnsmasq/g' /etc/NetworkManager/NetworkManager.conf || fatal "couldn't set dnsmasq as dns in NetworkManager";;
        *) sudo sed -i 's/^dns=.*$/dns=dnsmasq/g' /etc/NetworkManager/NetworkManager.conf || fatal "couldn't set dnsmasq as dns in NetworkManager";;
    esac
}

configure_start_dnsmasq () {
    sudo mkdir -p /etc/NetworkManager/dnsmasq.d/ || fatal "couln't create dnsmasq config directory"
    sudo bash -c 'echo "address=/.massa/127.0.0.1" > /etc/NetworkManager/dnsmasq.d/massa.conf'
    sudo mv /etc/resolv.conf /etc/resolv.conf_backup_thyra_install || fatal "couldn't make /etc/resolv.conf backup"
    sudo ln -s /var/run/NetworkManager/resolv.conf /etc/resolv.conf || fatal "couln't make link NetworkManager resolver to /etc/resolv.conf"
    sudo systemctl restart NetworkManager || fatal "dnsmasq service failed to restart"
}

set_local_dns () {
    case $(sudo lsof -i :53 | sed -n 2p | sed 's/[[:space:]].*$//') in
        "")         (configure_network_manager || fatal "couldn't set dnsmasq as dns") && configure_start_dnsmasq || exit -1 ;;
        dnsmasq)    configure_start_dnsmasq || exit -1 ;;
        systemd-r)  warn "Your computer has systemd-resolver as a DNS resolver. Thyra needs dnsmasq to redirect .massa website." && \
                    read -p "Do you agree to install dnsmasq in place of systemd-resolver ? [y/n]" yn
                    if [ "$yn" != "y" ] && [ "$yn" != "yes" ]; then
                        fatal "Aborting."
                    fi
                    (configure_network_manager || fatal "couldn't set dnsmasq as dns") && configure_start_dnsmasq || exit -1 ;;
        *)          fatal "local DNS application unsupported." ;;
    esac
}

generate_certificate () {
    green "INFO" "Installing MKcert and generating HTTPS certificates:"
    curl -sL $MKCERT_URL -o mkcert
    chmod +x mkcert
    ./mkcert -install
    mkdir -p $THYRA_CONF_DIR/certs
    ./mkcert -cert-file $THYRA_CONF_DIR/certs/cert.pem -key-file $THYRA_CONF_DIR/certs/cert-key.pem my.massa
    rm mkcert
}

echo ""

green "INFO" "This installation script will install the last release of Thyra and will configure your local DNS to resolve .massa if needed."

check_supported_distro || exit 1

install_thyra || exit 1

generate_certificate || exit 1

ping -c 1 -t 1 test.massa &> /dev/null || set_local_dns || exit 1

green "SUCCESS" "Thyra is installed and the .massa TLD resolution is configured. You're free to go!!!"

echo ""
