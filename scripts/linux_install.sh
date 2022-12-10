#!/bin/bash

set -x

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

install_deps () {
    sudo apt-get update
    # Needed for mkcert
    sudo apt install -y libnss3-tools
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
    green "INFO" "Set dnsmasq as NetworkManager dns resolver:"
    NM_CONF=/etc/NetworkManager/NetworkManager.conf
    # Backup config file
    sudo cp $NM_CONF "$NM_CONF"_backup_thyra_install
    # set dnsmasq plugin
    if grep '^dns=' "$NM_CONF"; then
        sudo sed -i '/\[main\]/adns=dnsmasq' $NM_CONF
    else
        sudo sed -i 's/^dns=.*$/dns=dnsmasq/g' $NM_CONF
    fi

    sudo ln -s /var/run/NetworkManager/resolv.conf /etc/resolv.conf

}

configure_dnsmasq_plugin () {
    green "INFO" "Register .massa localhost TLD:"
    sudo mkdir -p /etc/NetworkManager/dnsmasq.d/
    sudo bash -c 'echo "address=/.massa/127.0.0.1" > /etc/NetworkManager/dnsmasq.d/massa.conf'
}

configure_dnsmasq () {
    green "INFO" "Register .massa localhost TLD:"
    sudo mkdir -p /etc/dnsmasq.d/
    sudo bash -c 'echo "address=/.massa/127.0.0.1" > /etc/dnsmasq.d/massa.conf'
}

set_local_dns () {

    #checks if NetworkManager is installed
    if command -v nmcli &> /dev/null; then
        # setup dnsmasq plugin for NetworkManager
        configure_network_manager || fatal "couldn't configure network manager dnsmasq plugin"
        configure_dnsmasq_plugin || fatal "couldn't configure dnsmasq plugin"
        sudo systemctl restart NetworkManager
    elif sudo systemctl is-active systemd-resolved
        # setup dnsmasq with systemd-resolved
        # set dns server to localhost
        sudo sed -i 's/^#DNS=$/DNS=127.0.0.1/g' /etc/systemd/resolved.conf
        sudo apt-get install -y dnsmasq
        # option to avoid conflict when listening to port 53
        sudo sed -i 's/^#bind-interfaces$/bind-interfaces/g' /etc/dnsmasq.conf
        configure_dnsmasq
        sudo systemctl restart dnsmasq
        sudo systemctl restart systemd-resolved
    else
        fatal "DNS resolver not supported. Please contact Thyra support team"
    fi
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

green "INFO" "This installation script will install the last release of Thyra and will configure your local DNS to resolve .massa if needed."

check_supported_distro || exit 1

install_deps || exit 1

install_thyra || exit 1

generate_certificate || exit 1

ping -c 1 -t 1 test.massa &> /dev/null || set_local_dns

green "SUCCESS" "Thyra is installed and the .massa TLD resolution is configured. You're free to go!!!"
