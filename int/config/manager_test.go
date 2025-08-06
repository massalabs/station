package config

import (
	"sync"
	"testing"
	"time"

	"github.com/massalabs/station/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTest initializes the test environment
func setupTest(t *testing.T) {
	t.Helper()
	// Reset the singleton instance for each test
	ResetConfigManager()

	// Initialize logger for tests
	err := logger.InitializeGlobal("")
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
}

// TestMSConfigManager_Singleton tests the singleton pattern
func TestMSConfigManager_Singleton(t *testing.T) {
	setupTest(t)

	// Create test manager directly to avoid network calls
	manager := createTestConfigManager(t)

	// Manually set the singleton instance for testing
	instance = manager
	instanceOnce = sync.Once{} // Reset once to ensure it's been called
	instanceOnce.Do(func() {}) // Mark as called

	// Get manager multiple times
	manager1, err1 := GetConfigManager()
	require.NoError(t, err1)

	manager2, err2 := GetConfigManager()
	require.NoError(t, err2)

	manager3, err3 := NewMSConfigManager() // Test backward compatibility
	require.NoError(t, err3)

	// All should return the same instance
	assert.Same(t, manager1, manager2)
	assert.Same(t, manager1, manager3)
	assert.Same(t, manager2, manager3)
}

// TestMSConfigManager_ResetSingleton tests resetting the singleton
func TestMSConfigManager_ResetSingleton(t *testing.T) {
	setupTest(t)

	// Create test manager
	manager := createTestConfigManager(t)
	instance = manager
	instanceOnce = sync.Once{}
	instanceOnce.Do(func() {})

	// Get the instance
	manager1, err := GetConfigManager()
	require.NoError(t, err)
	assert.Same(t, manager, manager1)

	// Reset and verify it creates a new instance structure
	ResetConfigManager()

	// The instance should be reset
	instanceMutex.RLock()
	assert.Nil(t, instance)
	instanceMutex.RUnlock()
}

// TestMSConfigManager_CurrentNetwork tests getting the current network
func TestMSConfigManager_CurrentNetwork(t *testing.T) {
	setupTest(t)
	manager := createTestConfigManager(t)

	current := manager.CurrentNetwork()
	require.NotNil(t, current)
	assert.NotEmpty(t, current.Name)
	assert.NotEmpty(t, current.NodeURL)
}

// TestMSConfigManager_SwitchNetwork tests switching between networks
func TestMSConfigManager_SwitchNetwork(t *testing.T) {
	setupTest(t)
	manager := createTestConfigManager(t)

	// Get initial network
	initialNetwork := manager.CurrentNetwork()
	require.NotNil(t, initialNetwork)

	// Find a different network to switch to
	var targetNetworkName string
	for _, network := range manager.Network.Networks {
		if network.Name != initialNetwork.Name {
			targetNetworkName = network.Name
			break
		}
	}
	require.NotEmpty(t, targetNetworkName, "Need at least 2 networks for this test")

	// Switch to target network
	err := manager.SwitchNetwork(targetNetworkName)
	require.NoError(t, err)

	// Verify the switch
	currentNetwork := manager.CurrentNetwork()
	assert.Equal(t, targetNetworkName, currentNetwork.Name)
}

// TestMSConfigManager_SwitchNetwork_UnknownNetwork tests switching to an unknown network
func TestMSConfigManager_SwitchNetwork_UnknownNetwork(t *testing.T) {
	setupTest(t)
	manager := createTestConfigManager(t)

	// Try to switch to a non-existent network
	err := manager.SwitchNetwork("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown network")
}

// TestMSConfigManager_ConcurrentAccess tests concurrent access to the config manager
func TestMSConfigManager_ConcurrentAccess(t *testing.T) {
	setupTest(t)
	manager := createTestConfigManager(t)

	// Get available networks
	var networkNames []string
	for _, network := range manager.Network.Networks {
		networkNames = append(networkNames, network.Name)
	}
	require.Len(t, networkNames, 2, "Need at least 2 networks for concurrent test")

	// Run concurrent operations
	done := make(chan bool, 10)

	// Concurrent switches
	for i := 0; i < 5; i++ {
		go func(i int) {
			defer func() { done <- true }()
			networkName := networkNames[i%len(networkNames)]
			err := manager.SwitchNetwork(networkName)
			assert.NoError(t, err)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 5; i++ {
		go func() {
			defer func() { done <- true }()
			current := manager.CurrentNetwork()
			assert.NotNil(t, current)
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent operations")
		}
	}
}

// TestNetworkConfig_Structure tests the NetworkConfig structure
func TestNetworkConfig_Structure(t *testing.T) {
	setupTest(t)
	manager := createTestConfigManager(t)

	// Test that currentNetwork is properly set
	assert.NotNil(t, manager.Network.currentNetwork)

	// Test that Networks slice is populated
	assert.NotEmpty(t, manager.Network.Networks)

	// Test that each network has required fields
	for _, network := range manager.Network.Networks {
		assert.NotEmpty(t, network.Name)
		assert.NotEmpty(t, network.NodeURL)
		assert.NotZero(t, network.ChainID)
		// Status can be up or down, but should be set
		assert.True(t, network.status == NetworkStatusUp || network.status == NetworkStatusDown)
	}
}

// TestRPCInfos_Fields tests the RPCInfos structure
func TestRPCInfos_Fields(t *testing.T) {
	rpcInfo := RPCInfos{
		Name:    "test-network",
		NodeURL: "https://test.example.com",
		Version: "1.0.0",
		ChainID: 12345,
		status:  NetworkStatusUp,
	}

	assert.Equal(t, "test-network", rpcInfo.Name)
	assert.Equal(t, "https://test.example.com", rpcInfo.NodeURL)
	assert.Equal(t, "1.0.0", rpcInfo.Version)
	assert.Equal(t, uint64(12345), rpcInfo.ChainID)
	assert.Equal(t, NetworkStatusUp, rpcInfo.status)
}

// TestNetworkStatus_Constants tests the NetworkStatus constants
func TestNetworkStatus_Constants(t *testing.T) {
	assert.Equal(t, NetworkStatus("up"), NetworkStatusUp)
	assert.Equal(t, NetworkStatus("down"), NetworkStatusDown)
}

// createTestConfigManager creates a test config manager with mocked dependencies
func createTestConfigManager(t *testing.T) *MSConfigManager {
	// Create test networks
	networks := []RPCInfos{
		{
			Name:    "mainnet",
			NodeURL: "https://mainnet.massa.net/api/v2",
			Version: "1.0.0",
			ChainID: 77658377,
			status:  NetworkStatusUp,
		},
		{
			Name:    "buildnet",
			NodeURL: "https://buildnet.massa.net/api/v2",
			Version: "1.0.0",
			ChainID: 77658366,
			status:  NetworkStatusUp,
		},
	}

	manager := &MSConfigManager{
		Network: NetworkConfig{
			currentNetwork: &networks[0], // Default to mainnet
			Networks:       networks,
		},
	}

	return manager
}

// TestLoadConfig_DefaultConfig tests loading default configuration
func TestLoadConfig_DefaultConfig(t *testing.T) {
	// Skip this test as it requires file system operations
	// Use createTestConfigManager for unit testing instead
	t.Skip("Skipping integration test - use createTestConfigManager for unit tests")
}

// TestConfigFile_Structure tests the ConfigFile structure
func TestConfigFile_Structure(t *testing.T) {
	configFile := ConfigFile{
		Networks: map[string]RPCConfigItem{
			"test": {
				URL:     "https://test.example.com",
				Default: boolPtr(true),
			},
		},
	}

	assert.Contains(t, configFile.Networks, "test")
	testNetwork := configFile.Networks["test"]
	assert.Equal(t, "https://test.example.com", testNetwork.URL)
	require.NotNil(t, testNetwork.Default)
	assert.True(t, *testNetwork.Default)
}

// TestBoolPtr tests the boolPtr helper function
func TestBoolPtr(t *testing.T) {
	truePtr := boolPtr(true)
	falsePtr := boolPtr(false)

	require.NotNil(t, truePtr)
	require.NotNil(t, falsePtr)
	assert.True(t, *truePtr)
	assert.False(t, *falsePtr)
}
