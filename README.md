# thyra

[![CI](https://github.com/massalabs/thyra/actions/workflows/CI.yml/badge.svg)](https://github.com/massalabs/thyra/actions/workflows/CI.yml)
[![codecov](https://codecov.io/gh/massalabs/thyra/branch/main/graph/badge.svg?token=592LPZLC4M)](https://codecov.io/gh/massalabs/thyra)
[![Go Report Card](https://goreportcard.com/badge/github.com/massalabs/thyra)](https://goreportcard.com/report/github.com/massalabs/thyra)

An entrance to the Massa blockchain.


## Contribute
To contribute to this project, please read out please read our [Contributor's Guide](./CONTRIBUTING.md).

## Install thyra 
To install Thyra, please follow the instructions available in the [Installation Guide](./INSTALL.md).

## Install a plugin

Soon

-----

**Troubleshooting** If you got any problem at any steps of the process, join us on [Discord dedicated channel](https://discord.com/channels/828270821042159636/851942484212318259) or [report a problem](https://github.com/massalabs/thyra/issues/new/choose)

-----


## Going further 

### To pass options to `thyra-server`?

Thyra accepts different options that you can specify when you start the program.
In this section you will find a non-exhaustive list of such options and examples of how you can use them.

--node-server : Specify which Massa network Thyra will communicate with while running.
Accepts a URL, an IP address or one of the following values :

- TESTNET : Uses Massa's testnet
- LABNET : Uses Massa's labnet
- INNONET : Uses Massa's innonet
- LOCALHOST : Expect Massa's network to be hosted at 127.0.0.1

To use this option with a constant, you have to execute :
`thyra-server --node-server=LABNET`
To use this option with a custom IP address, you have to execute :
`thyra-server --node-server=192.168.X.X`

### To code with auto-reload

You can run the application with this command: `air`.

It will generate and reload thyra each time a file is modified.


## Additional information

### Why this name?

θύρα (thýra) in ancient Greek means door, entrance. This is exactly what this project is: an entrance to the Massa blockchain.

### How to pronounce it?

See <https://www.youtube.com/watch?v=_0BQ7sSJMTw>.
