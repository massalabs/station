# thyra

An entrance to the Massa blockchain.

## /!\ WIP

Everything is this project is WIP prototype.

/!\ Breaking changes ahead /!\

## Contribute

### Install dev dependencies

To develop on this project you will need :
- [go](https://go.dev/doc/install)
- [textFileToGoConst](https://github.com/logrusorgru/textFileToGoConst) to generate go constants from file contents
- [swagger](https://github.com/go-swagger/go-swagger) to generate go code from API documentation

Once Golang is installed on your system, you can install the last two dependencies by running the following command outside of a go module directory:
- `go install github.com/logrusorgru/textFileToGoConst@latest`
- `go install github.com/go-swagger/go-swagger/cmd/swagger@latest`

## How to ...

### ... get the last version of `thyra-server` ?

To install the last version of the `thyra-server` application you need to:

- [install go](https://go.dev/doc/install)
- execute `go install github.com/massalabs/thyra/cmd/thyra-server@main` in your terminal

Note: you can change `main` to a tag or a commit value if needed.

That's it, thyra-server is installed in '$HOME/go/bin/' dir and you can use it by executing `thyra-server` in your terminal.

### ... access a website stored on the Massa blockchain ?

Prerequisite: Having a running thyra-server application on your machine.

To access the website you need to go to http://localhost:8080/website?url=<address of the website>.
For instance, to access flappy text stored on the blockchain, click the following link: http://localhost:8080/website?url=A1aMywGBgBywiL6WcbKR4ugxoBtdP9P3waBVi5e713uvj7F1DJL.

### ... redirect massa Top Level Domain to localhost ?

#### Linux

1- Install dnsmasq

```shell
sudo apt install dnsmasq
```

2 - Add massa TLD resolution to localhost
Edit `/etc/dnsmasq.conf` and add `address=/.massa/127.0.0.1`

NOTE : If DNR is globally slow, add the following lines to the same file (`/etc/dnsmasq.conf`):

```shell
no-resolv
server=8.8.8.8
server=8.8.4.4
```

#### MacOS

See: https://www.larry.dev/no-more-etc-hosts-on-mac-with-dnsmasq/

#### Windows

As dnsmasq is not supported on windows, you can use Acrylic.

See: https://serverfault.com/questions/539591/how-to-resolve-all-dev-domains-to-localhost-on-windows#answer-808963

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

## Additional information

### Why this name ?

θύρα (thýra) in ancient Greek means door, entrance. This is exactly what this project is: an entrance to the Massa blockchain.

### How to pronounce it ?

See https://www.youtube.com/watch?v=_0BQ7sSJMTw.
