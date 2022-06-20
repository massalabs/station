# thyra
An entrance to the Massa blockchain.

## /!\ WIP

Everything is this project is WIP prototype.

/!\ Breaking changes ahead /!\

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




## Additional information

### Why this name ?

θύρα (thýra) in ancient Greek means door, entrance. This is exactly what this project is: an entrance to the Massa blockchain.

### How to pronounce it ?

See https://www.youtube.com/watch?v=_0BQ7sSJMTw.
