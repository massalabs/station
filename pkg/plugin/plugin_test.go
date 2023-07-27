package plugin

import (
	"testing"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/store"
	"github.com/stretchr/testify/assert"
)

type testPluginUpdate struct {
	plugin              Plugin
	storePlugins        *store.Store
	updatable           bool
	massaStationVersion string
	compatible          bool
}

//nolint:govet
func AssertUpdate(t *testing.T, testPluginUpdate testPluginUpdate) {
	t.Helper()

	storeMS := testPluginUpdate.storePlugins
	pluginInStore := storeMS.Plugins[0]
	plgn := &testPluginUpdate.plugin
	plugins := make(map[string]*Plugin)

	plugins[plgn.ID] = plgn
	plgn.status = Up

	config.Version = testPluginUpdate.massaStationVersion

	pluginCompatible, err := pluginInStore.IsPluginCompatible()
	assert.NoError(t, err)

	pluginInStore.IsCompatible = pluginCompatible
	assert.Equal(t, pluginCompatible, testPluginUpdate.compatible)

	storeMS.Plugins[0].IsCompatible = pluginCompatible
	isUpdatable, err := storeMS.CheckForPluginUpdates(plgn.info.Name, plgn.info.Version)
	assert.NoError(t, err)
	assert.Equal(t, isUpdatable, testPluginUpdate.updatable)
}

//nolint:funlen
func TestPlugin_Update(t *testing.T) {
	//nolint:exhaustruct
	pluginNonUpdatable := Plugin{
		ID: "test",
		info: &Information{
			Name:    "test",
			Version: "10.0.0",
		},
	}
	//nolint:exhaustruct
	pluginUpdatable := Plugin{
		ID: "test",
		info: &Information{
			Name:    "test",
			Version: "0.0.19",
		},
	}
	baseStore := &store.Store{
		Plugins: []store.Plugin{
			{
				Name:                "test",
				Version:             "1.0.9",
				MassaStationVersion: ">=1.0.0",
			},
		},
	}
	//nolint:govet
	testsPluginUpdate := []testPluginUpdate{
		{
			plugin:              pluginNonUpdatable,
			storePlugins:        baseStore,
			updatable:           false,
			massaStationVersion: "1.0.0",
			compatible:          true,
		},
		//nolint:govet

		{
			plugin:              pluginNonUpdatable,
			storePlugins:        baseStore,
			updatable:           false,
			massaStationVersion: "0.1.0",
			compatible:          false,
		},
		//nolint:govet

		{
			plugin:              pluginUpdatable,
			storePlugins:        baseStore,
			updatable:           true,
			massaStationVersion: "1.1.0",
			compatible:          true,
		},
		//nolint:govet

		{
			plugin:              pluginUpdatable,
			storePlugins:        baseStore,
			updatable:           true,
			massaStationVersion: "dev", // dev version is always compatible
			compatible:          true,
		},
	}

	//nolint:govet
	for _, test := range testsPluginUpdate {
		AssertUpdate(t, test)
	}
}
