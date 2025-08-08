package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setTempHome configures HOME and XDG_CONFIG_HOME to a temporary directory for tests
func setTempHome(t *testing.T) (restore func()) {
	t.Helper()
	tempDir := t.TempDir()

	oldHome := os.Getenv("HOME")
	oldXDG := os.Getenv("XDG_CONFIG_HOME")

	require.NoError(t, os.Setenv("HOME", tempDir))
	require.NoError(t, os.Setenv("XDG_CONFIG_HOME", tempDir))

	return func() {
		_ = os.Setenv("HOME", oldHome)
		_ = os.Setenv("XDG_CONFIG_HOME", oldXDG)
	}
}

// newTestManagerFromDefaults builds an in-memory manager seeded from DefaultConfig without network calls
func newTestManagerFromDefaults(t *testing.T) *MSConfigManager {
	t.Helper()
	networks := make([]RPCInfos, 0, len(DefaultConfig.Networks))
	var current RPCInfos
	for name, item := range DefaultConfig.Networks {
		net := RPCInfos{
			Name:    name,
			NodeURL: item.URL,
			Version: "", // not needed for tests
			ChainID: 0,  // not needed for tests
			status:  NetworkStatusUp,
		}
		networks = append(networks, net)
		if item.Default != nil && *item.Default {
			current = net
		}
	}
	// Fallback: if no default flagged in DefaultConfig, take first as current
	if current.Name == "" && len(networks) > 0 {
		current = networks[0]
	}

	return &MSConfigManager{
		Network: NetworkConfig{
			currentNetwork: &current,
			Networks:       networks,
		},
	}
}

func TestAddNetwork_PersistsAndSetsDefault(t *testing.T) {
	setupTest(t)
	restore := setTempHome(t)
	defer restore()

	mgr := newTestManagerFromDefaults(t)

	// Add a new custom network and make it default
	err := mgr.AddNetwork("custom", "http://localhost:12345/api", true)
	require.NoError(t, err)

	// In-memory assertions
	var found bool
	for _, n := range mgr.Network.Networks {
		if n.Name == "custom" {
			found = true
			break
		}
	}
	assert.True(t, found, "custom network should be present in memory")
	require.NotNil(t, mgr.CurrentNetwork())
	assert.Equal(t, "custom", mgr.CurrentNetwork().Name)

	// Persisted config assertions
	cfg, err := LoadConfig()
	require.NoError(t, err)
	item, ok := cfg.Networks["custom"]
	if assert.True(t, ok, "custom network should be persisted") {
		assert.Equal(t, "http://localhost:12345/api", item.URL)
		require.NotNil(t, item.Default)
		assert.True(t, *item.Default)
	}
	// Others should not be default
	for name, v := range cfg.Networks {
		if name == "custom" {
			continue
		}
		if v.Default != nil {
			assert.False(t, *v.Default, "only custom should be default")
		}
	}
}

func TestEditNetwork_RenameURLAndDefault(t *testing.T) {
	setupTest(t)
	restore := setTempHome(t)
	defer restore()

	mgr := newTestManagerFromDefaults(t)

	// Ensure the default config file exists
	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotEmpty(t, cfg.Networks)

	// Choose an existing network to edit (e.g., buildnet if present)
	target := ""
	for name := range cfg.Networks {
		if name != "mainnet" {
			target = name
			break
		}
	}
	if target == "" { // fallback to any
		for name := range cfg.Networks {
			target = name
			break
		}
	}
	require.NotEmpty(t, target)

	newURL := "http://localhost:18000/api"
	newName := "devnet"
	makeDefault := true

	err = mgr.EditNetwork(target, &newURL, &makeDefault, &newName)
	require.NoError(t, err)

	// In-memory assertions
	var renamedFound bool
	for _, n := range mgr.Network.Networks {
		if n.Name == newName {
			renamedFound = true
			assert.Equal(t, newURL, n.NodeURL)
		}
		assert.NotEqual(t, target, n.Name)
	}
	assert.True(t, renamedFound, "renamed network should be present in memory")
	require.NotNil(t, mgr.CurrentNetwork())
	assert.Equal(t, newName, mgr.CurrentNetwork().Name)

	// Persisted config assertions
	cfg, err = LoadConfig()
	require.NoError(t, err)
	_, okOld := cfg.Networks[target]
	assert.False(t, okOld, "old name should be removed from persisted config")
	newItem, okNew := cfg.Networks[newName]
	if assert.True(t, okNew, "new name should be present in persisted config") {
		assert.Equal(t, newURL, newItem.URL)
		require.NotNil(t, newItem.Default)
		assert.True(t, *newItem.Default)
	}
}

func TestDeleteNetwork_RemovesAndSwitchesCurrent(t *testing.T) {
	setupTest(t)
	restore := setTempHome(t)
	defer restore()

	mgr := newTestManagerFromDefaults(t)

	// Persisted config should exist at test location
	cfgPath, err := getConfigFilePath()
	require.NoError(t, err)
	info, err := os.Stat(filepath.Dir(cfgPath))
	require.NoError(t, err)
	require.True(t, info.IsDir())

	// Ensure we have at least 2 networks
	require.GreaterOrEqual(t, len(mgr.Network.Networks), 2, "test requires at least 2 networks")
	toDelete := mgr.Network.Networks[1].Name

	// Set current to the one we will delete to test switching
	mgr.Network.currentNetwork = &mgr.Network.Networks[1]

	err = mgr.DeleteNetwork(toDelete)
	require.NoError(t, err)

	// In-memory: deleted, and current switched
	for _, n := range mgr.Network.Networks {
		assert.NotEqual(t, toDelete, n.Name)
	}
	require.NotNil(t, mgr.CurrentNetwork())
	assert.NotEqual(t, toDelete, mgr.CurrentNetwork().Name)

	// Persisted: deleted and there is at least one default
	cfg, err := LoadConfig()
	require.NoError(t, err)
	_, ok := cfg.Networks[toDelete]
	assert.False(t, ok, "deleted network should not be in persisted config")

	hasDefault := false
	for _, v := range cfg.Networks {
		if v.Default != nil && *v.Default {
			hasDefault = true
			break
		}
	}
	assert.True(t, hasDefault, "one remaining network should be default")
}

func TestDeleteNetwork_LastRemainingFails(t *testing.T) {
	setupTest(t)
	restore := setTempHome(t)
	defer restore()

	// Persist config with only one network
	cfg, err := LoadConfig()
	require.NoError(t, err)
	cfg.Networks = map[string]RPCConfigItem{
		"onlyone": {URL: "http://localhost:19000/api", Default: boolPtr(true)},
	}
	require.NoError(t, saveConfig(cfg))

	// In-memory manager with a single network
	mgr := &MSConfigManager{
		Network: NetworkConfig{
			currentNetwork: &RPCInfos{Name: "onlyone", NodeURL: "http://localhost:19000/api", status: NetworkStatusUp},
			Networks:       []RPCInfos{{Name: "onlyone", NodeURL: "http://localhost:19000/api", status: NetworkStatusUp}},
		},
	}

	err = mgr.DeleteNetwork("onlyone")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete the last remaining network")
}
