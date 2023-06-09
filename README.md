# MassaStation: An entrance to the Massa blockchain

[![CI](https://github.com/massalabs/thyra/actions/workflows/api.yml/badge.svg)](https://github.com/massalabs/thyra/actions/workflows/api.yml)
[![codecov](https://codecov.io/gh/massalabs/thyra/branch/main/graph/badge.svg?token=592LPZLC4M)](https://codecov.io/gh/massalabs/thyra)
[![Go Report Card](https://goreportcard.com/badge/github.com/massalabs/thyra)](https://goreportcard.com/report/github.com/massalabs/thyra)


# Install MassaStation

To install MassaStation, please follow the instructions available in the [Installation Guide](./INSTALLATION.md).


# Contributing

We welcome contributions of all kinds, from bug reports to feature requests and code contributions.

If you're interested in contributing to MassaStation, please make sure to read our [Contribution Guidelines](./CONTRIBUTING.md) for detailed instructions on how to get started. 


# Going further 

## Install plugins

MassaStation can be extended with plugins. Plugins that were validated by Massa Labs are available in the [Plugin Store](https://station.massa/store/).

If you want to install a plugin that is not available in the Plugin Store, you can install it manually:
1. Get the `.zip` file download URL of the plugin you want to install. Make sure this URL matches the version of MassaStation you are using, your computer OS and architecture.
2. Paste the URL in the `Install a plugin` field of the [plugin manager page](https://station.massa/store/).
3. Click on the `Install` button.

The plugin will be installed and activated after a few seconds. You can now use it in MassaStation !


## Select the network to use

MassaStation can be configured to use your own node or one of Massa's networks. To do so, you have to use the `--node-server` option. This option accepts a URL, an IP address or one of the following values :

- `TESTNET` : Connects MassaStation to a node running on Massa's testnet. This is the default value.
- `LABNET` : Connects MassaStation to a node running on Massa's labnet.
- `BUILDNET` : Connects MassaStation to a node running on Massa's buildnet.
- `LOCALHOST` : Connects MassaStation to the node running on the same machine as MassaStation, and so, available at `127.0.0.1`.

It is also possible to use a custom IP address or URL:
`massastation --node-server={IP_ADDRESS}`
