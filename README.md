# thyra

[![CI](https://github.com/massalabs/thyra/actions/workflows/CI.yml/badge.svg)](https://github.com/massalabs/thyra/actions/workflows/CI.yml)
[![codecov](https://codecov.io/gh/massalabs/thyra/branch/main/graph/badge.svg?token=592LPZLC4M)](https://codecov.io/gh/massalabs/thyra)
[![Go Report Card](https://goreportcard.com/badge/github.com/massalabs/thyra)](https://goreportcard.com/report/github.com/massalabs/thyra)

An entrance to the Massa blockchain.


## Contribute
To contribute to this project, please read our [Contributor's Guide](./CONTRIBUTING.md).


## Install thyra 
To install Thyra, please follow the instructions available in the [Installation Guide](./INSTALLATION.md).

## Install a plugin

Copy/paste the latest release .zip file to Thyra [plugin manager page](https://my.massa/thyra/plugin-manager/)

## Going further 

### Use your own node with Thyra

Thyra accepts different options that you can specify when you start the program. Like this you will be able to use your own node, for instance.
In this section you will find a non-exhaustive list of options and examples of how you can use them.


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


## Additional information

### Why this name?

θύρα (thýra) in ancient Greek means door, entrance. This is exactly what this project is: an entrance to the Massa blockchain.

### How to pronounce it?

See <https://www.youtube.com/watch?v=_0BQ7sSJMTw>.
