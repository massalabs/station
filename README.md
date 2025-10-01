# MassaStation: An entrance to the Massa blockchain

[![CI](https://github.com/massalabs/station/actions/workflows/api.yml/badge.svg?branch=main)](https://github.com/massalabs/station/actions/workflows/api.yml?query=branch%3Amain)
[![codecov](https://codecov.io/gh/massalabs/station/branch/main/graph/badge.svg?token=592LPZLC4M)](https://codecov.io/gh/massalabs/station)
[![Go Report Card](https://goreportcard.com/badge/github.com/massalabs/station)](https://goreportcard.com/report/github.com/massalabs/station)

## Install MassaStation

To install MassaStation, please follow the instructions available in the [Installation Guide](https://docs.massa.net/docs/massastation/install).

## Contributing

We welcome contributions of all kinds, from bug reports to feature requests and code contributions.

If you're interested in contributing to MassaStation, please make sure to read our [Contribution Guidelines](./CONTRIBUTING.md) for detailed instructions on how to get started.

## Going further

### Plugins

MassaStation is a plugin manager. It enables everyone to use, create and enjoy different plugins to activate features to the Massa blockchain.

#### Install a plugin

You can install plugins that were validated by Massa Labs from the [Plugin Store](https://station.massa/web/store).

The plugin will be automatically installed and activated after a few seconds directly in your Station. Browse MassaStation store to find the plugin you need.

#### Create a Plugin

##### development
If you are working on a plugin, you can use logic from [plugin-kit](./plugin-kit/) go module.  
It provides code logic for your plugin to register it-self to Massa Station plugin manager.
Here is how to use it in your plugin:

```golang
	pluginKit "github.com/massalabs/station/plugin-kit"
    ...

    listener, err := server.HTTPListener()
	if err != nil {
		panic(err)
	}

    pluginKit.RegisterPlugin(listener)
```

##### test
You can install your plugin manually to test it using MassaStation:

1. Get the `.zip` file download URL of the plugin you want to install. Make sure this URL matches the version of MassaStation you are using, your computer OS and architecture.
2. Paste the URL in the `Install a plugin` field of the [plugin manager page](https://station.massa/web/store).
3. Click on the `Install` button.


> **Note:** A complete guide on how to create a plugin is available [here](https://docs.massa.net/docs/massaStation/guidelines)



## License
Please make sure to read our [Software License](./LICENSE.md)
