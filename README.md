# thyra

[![CI](https://github.com/massalabs/thyra/actions/workflows/CI.yml/badge.svg)](https://github.com/massalabs/thyra/actions/workflows/CI.yml)
[![codecov](https://codecov.io/gh/massalabs/thyra/branch/main/graph/badge.svg?token=592LPZLC4M)](https://codecov.io/gh/massalabs/thyra)
[![Go Report Card](https://goreportcard.com/badge/github.com/massalabs/thyra)](https://goreportcard.com/report/github.com/massalabs/thyra)

An entrance to the Massa blockchain.

## ⚠️ WIP

This project is still WIP. It is a prototype. 

⚠️ Potential breaking changes ahead ⚠️

## Contribute

### Install dev dependencies

To develop on this project you will need :

- [go](https://go.dev/doc/install)
- [swagger](https://github.com/go-swagger/go-swagger) to generate go code from API documentation

Once Golang is installed on your system, you can install the swagger dependency by running the following command outside of a go module directory:

- `go install github.com/go-swagger/go-swagger/cmd/swagger@latest`

### Setup code formatting tool

Use prettrier to format de code. We recommend to install IDE prettier extension to format on save.

For go code, we use the formatter included in <https://marketplace.visualstudio.com/items?itemName=golang.go>.


### golangci-lint

- golangci-lint is used to run linters in parallel.. We recommend to [install](https://golangci-lint.run/usage/install/) it locally and run it on your source code, before pushing any modification, otherwise some potential lint errors will be catched by the pipeline.
- to run golangci-lint locally : `golangci-lint run .`

#### How to resolve golangci-lint recurring errors ?

- File is not `gofumpt`

gofumpt need to be installed locally `go install mvdan.cc/gofumpt@latest`

run gofumpt locally on your source code `gofumpt -l -w .`

- File is not `gci`

gci need to be installed locally `go install github.com/daixiang0/gci@latest`

run gofumpt locally on your source code `gci --write .`

## How to...

### ... install Thyra on my computer?

Follow the instructions for your computer in the wiki:

- [MacOS](https://github.com/massalabs/thyra/blob/main/INSTALLATION.md#macos)
- [Linux](https://github.com/massalabs/thyra/blob/main/INSTALLATION.md#linux)
- [Windows](https://github.com/massalabs/thyra/blob/main/INSTALLATION.md#windows)

### ... manage my wallet?
1. Create / delete your wallet 

You can access to Thyra wallet interface at URL : <http://my.massa/thyra/wallet/index.html>
By inputing the 'Nickname' & 'Password', you'll be able to create an encrypted wallet locally on your machine.
To delete your wallet, simply use the interface. 
⚠️ If you delete your wallet, you won't be able to edit the website linked to it anymore.

2. Get coins on your wallet

To get coins on your wallet, you have to send your address on [Massa faucet channel](https://discord.com/channels/828270821042159636/866190913030193172)
Make sure that you use the latest version of Thyra (and defacto Testnet) [here](https://github.com/massalabs/thyra/releases/latest/), otherwise the faucet won't work.

### ... store a website on chain?

You can access to Thyra web hosting interface at URL : <http://my.massa/thyra/websiteCreator/index.html>

In order to register a website on Thyra you'll need to :

- Deploy a Smart Contract that will handle the storage of your website, your DNS name will fetch the Address of this Smart Contract
- Upload the build of your application
- Use a wallet with sufficient coins to upload it on the blockchain
Important note: At the moment, we have defined that 1 chunk (=280ko) of data worth 100 MAS. It will change and become more and more specific and precise as the Testnet and Thyra are evolving. In the mean time, we have defined it arbitrarily.
- Share your .massa websites on our [Discord channel](https://discord.com/channels/828270821042159636/912346860902047755) !

### ... get the latest dev version of `thyra-server`?

To install the latest dev version of the `thyra-server` application you need to:

- [install go](https://go.dev/doc/install)
- execute `go install github.com/massalabs/thyra/cmd/thyra-server@main` in your terminal

Note: you can change `main` to a tag or a commit value if needed.

That's it, thyra-server is installed in '$HOME/go/bin/' dir and you can use it by executing `thyra-server` in your terminal.

### ... pass options to `thyra-server`?

Thyra accepts different options that you can specify when you start the program.
In this section you will find a non-exhaustive list of such options and examples of how you can use them.

--node-server : Specify which Massa network Thyra will communicate with while running.
Accepts a URL, an IP address or one of the following values :

- TESTNET : Uses Massa's testnet
- LABNET : Uses Massa's labnet
- INNONET : Uses Massa's innonet
- LOCALHOST : Expect Massa's network to be hosted at 127.0.0.1

To use this option with a constant, you have to execute :
`thyra-server node-server=LABNET`
To use this option with a custom IP address, you have to execute :
`thyra-server node-server=192.168.X.X`

### ... secure HTTPS configuration?

Using HTTPS configuration without specifying your own certificate and key triggers a warning: `insecure HTTPS configuration`.

To solve this you need to create your own certificate. You can do so by using openssl:

```shell
openssl req -newkey rsa:4096 \
            -x509 \
            -sha256 \
            -days 365 \
            -nodes \
            -out my_thyra.crt \
            -keyout my_thyra.key
```

You can now execute a thyra-server using the following command changing _path to ..._ to proper values:
`thyra-server --tls-certificate <path to my_thyra.crt> --tls-key <path to my_thyra.key>`.

### ... code with auto-reload

You can run the application with this command: `air`.

It will generate and reload thyra each time a file is modified.

## Usage

### Upload a website

With Thyra you can upload a website. It will be hosted on the Massa blockchain and available as a `.massa` URL. You can access this URL through Tyhra.

For this to work, the file you upload must be a zip archive (file ending with `.zip`). This archive must contain a `index.html` file at the root.

## Additional information

### Why this name?

θύρα (thýra) in ancient Greek means door, entrance. This is exactly what this project is: an entrance to the Massa blockchain.

### How to pronounce it?

See <https://www.youtube.com/watch?v=_0BQ7sSJMTw>.
