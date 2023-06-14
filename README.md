# MassaStation: An entrance to the Massa blockchain

[![CI](https://github.com/massalabs/thyra/actions/workflows/api.yml/badge.svg?branch=main)](https://github.com/massalabs/thyra/actions/workflows/api.yml?query=branch%3Amain)
[![codecov](https://codecov.io/gh/massalabs/thyra/branch/main/graph/badge.svg?token=592LPZLC4M)](https://codecov.io/gh/massalabs/thyra)
[![Go Report Card](https://goreportcard.com/badge/github.com/massalabs/thyra)](https://goreportcard.com/report/github.com/massalabs/thyra)


## Install MassaStation

To install MassaStation, please follow the instructions available in the [Installation Guide](./INSTALLATION.md).


## Contributing

We welcome contributions of all kinds, from bug reports to feature requests and code contributions.

If you're interested in contributing to MassaStation, please make sure to read our [Contribution Guidelines](./CONTRIBUTING.md) for detailed instructions on how to get started. 


## Going further

### Modules (plugins)

MassaStation is a module manager. It enables everyone to use, create and enjoy different modules to activate features to the Massa blockchain.

#### Install a module

You can install modules that were validated by Massa Labs from the [Module Store](https://station.massa/store/).

The module will be automatically installed and activated after a few seconds directly in your Station. Browse MassaStation store to find the module you need.


#### Create a module

If you are working on a module, you can install it manually to test it using MassaStation:
1. Get the `.zip` file download URL of the module you want to install. Make sure this URL matches the version of MassaStation you are using, your computer OS and architecture.
2. Paste the URL in the `Install a plugin` field of the [module manager page](https://station.massa/store/).
3. Click on the `Install` button.

> **Note:** A complete guide on how to create a module will be available soon.


### Network

MassaStation can be configured to use your own node or one of Massa's networks. To do so, you have to use the `--node-server` option. This option accepts a URL, an IP address or one of the following values :

- `TESTNET` : Connects MassaStation to a node running on Massa's testnet. This is the default value.
- `LABNET` : Connects MassaStation to a node running on Massa's labnet.
- `BUILDNET` : Connects MassaStation to a node running on Massa's buildnet.
- `LOCALHOST` : Connects MassaStation to the node running on the same machine as MassaStation, and so, available at `127.0.0.1`.

It is also possible to use a custom IP address or URL:
`massastation --node-server={IP_ADDRESS}`
