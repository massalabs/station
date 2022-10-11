# thyra

[![CI](https://github.com/massalabs/thyra/actions/workflows/CI.yml/badge.svg)](https://github.com/massalabs/thyra/actions/workflows/CI.yml)
[![codecov](https://codecov.io/gh/massalabs/thyra/branch/main/graph/badge.svg?token=592LPZLC4M)](https://codecov.io/gh/massalabs/thyra)
[![Go Report Card](https://goreportcard.com/badge/github.com/massalabs/thyra)](https://goreportcard.com/report/github.com/massalabs/thyra)

An entrance to the Massa blockchain.

## /!\ WIP

Everything is this project is WIP prototype.

/!\ Breaking changes ahead /!\

## Contribute

### Install dev dependencies

To develop on this project you will need :

-   [go](https://go.dev/doc/install)
-   [swagger](https://github.com/go-swagger/go-swagger) to generate go code from API documentation

Once Golang is installed on your system, you can install the swagger dependency by running the following command outside of a go module directory:

-   `go install github.com/go-swagger/go-swagger/cmd/swagger@latest`

## How to ...

### ... install Thyra on my computer ?

Follow the instructions for your computer in the wiki:

-   [MacOS](https://github.com/massalabs/thyra/blob/main/INSTALLATION.md#macos)
-   [Linux](https://github.com/massalabs/thyra/blob/main/INSTALLATION.md#linux)
-   [Windows](https://github.com/massalabs/thyra/blob/main/INSTALLATION.md#windows)

### ... get the latest dev version of `thyra-server` ?

To install the latest dev version of the `thyra-server` application you need to:

-   [install go](https://go.dev/doc/install)
-   execute `go install github.com/massalabs/thyra/cmd/thyra-server@main` in your terminal

Note: you can change `main` to a tag or a commit value if needed.

That's it, thyra-server is installed in '$HOME/go/bin/' dir and you can use it by executing `thyra-server` in your terminal.

### ... pass options to `thyra-server` ?

Thyra accepts different options that you can specify when you start the program.
In this section you will find a non-exhaustive list of such options and examples of how you can use them.

--node-server : Specify which Massa network Thyra will communicate with while runnning.
Accepts a URL, an IP address or one of the following values :

-   TESTNET : Uses Massa's testnet
-   LABNET : Uses Massa's labnet
-   INNONET : Uses Massa's innonet
-   LOCALHOST : Expect Massa's network to be hosted at 127.0.0.1

To use this option with a constant, you have to execute :
`thyra-server node-server=LABNET`
To use this option with a custom IP address, you have to execute :
`thyra-server node-server=192.168.X.X`

### ... secure HTTPS configuration ?

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

### ... wallet management ?

You can access to Thyra wallet interface at URL : http://my.massa/thyra/wallet/index.html
By inputing the 'Nickname' & 'Password', you'll be able to create an encrypted wallet locally on your machine.

### ... web on chain ?

You can access to Thyra web hosting interface at URL : http://my.massa/thyra/websiteCreator/index.html

In order to register a website on Thyra you'll need to :

-   Deploy a Smart Contract that will handle the storage of your website, your DNS name will fetch the Address of this Smart Contract
-   Upload the build of your application

## Usage

### Upload a website

With Thyra you can upload a website. It will be hosted on the Massa blockchain and available as a `.massa` URL. You can access this URL through Tyhra.

For this to work, the file you upload must be a zip archive (file ending with `.zip`). This archive must contain a `index.html` file at the root.

## Additional information

### Why this name ?

θύρα (thýra) in ancient Greek means door, entrance. This is exactly what this project is: an entrance to the Massa blockchain.

### How to pronounce it ?

See https://www.youtube.com/watch?v=_0BQ7sSJMTw.
