#!/bin/bash

## Detect architecture
architecture=""
case $(uname -m) in
    x86_64) architecture="amd64" ;;
    arm64)  architecture="arm64" ;;
    *)      echo "Error: Unsupported architecture." && exit 1 ;;
esac

## Downloading Thyra
curl -L "https://github.com/massalabs/thyra/releases/latest/download/thyra-server_darwin_${architecture}" -o thyra-server
chmod +x ./thyra-server

## Install DNSMasq
brew install dnsmasq

sudo bash -c 'echo "address=/.massa/127.0.0.1" > $(brew --prefix)/etc/dnsmasq.d/massa.conf'
sudo mkdir -p /etc/resolver
sudo bash -c 'echo "nameserver 127.0.0.1" > /etc/resolver/massa'

sudo brew services start dnsmasq
