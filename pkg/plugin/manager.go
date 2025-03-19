package plugin

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin/utils"
	"github.com/massalabs/station/pkg/store"
)

// Directory returns the plugin directory.
// Note: the plugin directory is the /plugins inside the home directory.
func Directory(configDir string) string {
	pluginsDir := path.Join(configDir, "plugins")
	_, err := os.Stat(pluginsDir)

	if os.IsNotExist(err) {
		err := os.MkdirAll(pluginsDir, os.ModePerm)
		if err != nil {
			panic(fmt.Errorf("getting plugins directory: creating folder: %w", err))
		}
	}

	return pluginsDir
}

// Manager manages different plugins.
// plugins key is a plugin map storage using the author name and the plugin name as key.
// correlationID is an identifier used to recognize the plugin when it register.
type Manager struct {
	mutex          sync.RWMutex
	plugins        map[string]*Plugin
	authorNameToID map[string]string
	configDir      string
}

// NewManager instantiates a manager struct.
func NewManager(configDir string) *Manager {
	//nolint:exhaustruct
	manager := &Manager{
		plugins:        make(map[string]*Plugin),
		authorNameToID: make(map[string]string),
		configDir:      configDir,
	}

	return manager
}

// ID returns the list of all the plugin correlationID.
func (m *Manager) ID() []string {
	keys := make([]string, len(m.plugins))

	i := 0

	for k := range m.plugins {
		keys[i] = k
		i++
	}

	return keys
}

// SetAlias adds an alias to an existing plugin.
// Alias can be defined during plugin register once the name and author of the plugin can be found.
func (m *Manager) SetAlias(alias string, correlationID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.plugins[correlationID] == nil {
		return fmt.Errorf("while setting alias for %s: no plugin matching the given correlationID %s", alias, correlationID)
	}

	registeredID, exist := m.authorNameToID[alias]
	if exist && registeredID != correlationID {
		return fmt.Errorf("while setting alias for %s: a plugin with the same alias already exists", alias)
	}

	m.authorNameToID[alias] = correlationID

	return nil
}

func (m *Manager) RemoveAlias(alias string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.authorNameToID, alias)
}

// PluginByAlias returns a plugin from the manager using an alias.
func (m *Manager) PluginByAlias(alias string) (*Plugin, error) {
	correlationID, exist := m.authorNameToID[alias]
	if exist {
		p, err := m.Plugin(correlationID)
		if err != nil {
			return nil, fmt.Errorf("getting plugin by alias %w", err)
		}

		return p, nil
	}

	return nil, fmt.Errorf("plugin not found for alias %s", alias)
}

// Plugin returns a plugin from the manager.
func (m *Manager) Plugin(correlationID string) (*Plugin, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	p, ok := m.plugins[correlationID]
	if !ok {
		return nil, fmt.Errorf("no plugin matching correlationID %s", correlationID)
	}

	return p, nil
}

// Unregister a plugin from the manager.
func (m *Manager) UnregisterPlugin(correlationID string) error {
	_, ok := m.plugins[correlationID]
	if !ok {
		return fmt.Errorf("no plugin matching correlationID %s", correlationID)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.plugins, correlationID)

	return nil
}

// Delete kills a plugin and remove it from the manager.
func (m *Manager) Delete(correlationID string) error {
	plgn, err := m.Plugin(correlationID)
	if err != nil {
		return fmt.Errorf("getting plugin %s: %w", correlationID, err)
	}

	// Ignore Stop errors. We want to delete the plugin anyway
	err = m.StopPlugin(plgn, true)
	if err != nil {
		logger.Warnf("stopping plugin before delete %s: %s\n", correlationID, err)
	}

	err = removePlugin(plgn)
	if err != nil {
		return fmt.Errorf("deleting plugin %s: %w", correlationID, err)
	}

	return nil
}

// removePlugin removes all the plugin files from the file system.
func removePlugin(plugin *Plugin) error {
	dir := plugin.Path

	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("removing plugin dir %s: %w", dir, err)
	}

	return nil
}

// generateCorrelationID generate a unique correlation correlationID.
func (m *Manager) generateCorrelationID() string {
	for {
		var maxInt int64 = 10_000
		idInteger, _ := rand.Int(rand.Reader, big.NewInt(maxInt))

		correlationID := idInteger.String()

		_, exist := m.plugins[correlationID]
		if exist {
			continue
		}

		return correlationID
	}
}

// InitPlugin starts new plugin and adds it to manager.
func (m *Manager) InitPlugin(pluginPath string) error {
	correlationID := m.generateCorrelationID()

	plugin, err := New(pluginPath, correlationID)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	m.plugins[correlationID] = plugin
	m.mutex.Unlock()

	return nil
}

// RunAll runs all the installed plugins.
func (m *Manager) RunAll() error {
	pluginDir := Directory(m.configDir)

	rootItems, err := os.ReadDir(pluginDir)
	if err != nil {
		return fmt.Errorf("reading plugins directory '%s': %w", pluginDir, err)
	}

	for _, rootItem := range rootItems {
		if rootItem.IsDir() {
			pluginName := rootItem.Name()
			pluginPath := filepath.Join(pluginDir, pluginName)

			err = m.InitPlugin(pluginPath)
			if err != nil {
				logger.Warnf("While running plugin %s: %s.", pluginName, err)
				logger.Warnf("This plugin will not be executed.")
			}
		}
	}

	return nil
}

func (m *Manager) StopAll() {
	logger.Info("Stopping all plugins...")

	for _, plugin := range m.plugins {
		err := m.StopPlugin(plugin, true)
		if err != nil {
			logger.Warnf("stopping plugin %s: %w", plugin.ID, err)
		}
	}
}

func (m *Manager) StopPlugin(plugin *Plugin, unregister bool) error {
	err := plugin.Stop()
	if err != nil {
		msg := fmt.Sprintf("error while stopping plugin %s: %s", plugin.info.Name, err)
		// if we want to unregister, then this error is not blocking
		if unregister {
			logger.Info(msg)
		} else {
			return fmt.Errorf(msg)
		}
	}

	alias := Alias(plugin.info.Author, plugin.info.Name)

	m.RemoveAlias(alias)

	if unregister {
		err = m.UnregisterPlugin(plugin.ID)
		if err != nil {
			return fmt.Errorf("unregistering plugin %s: %w", plugin.ID, err)
		}
	}

	return nil
}

// Install grabs a remote plugin from the given url and install it locally.
func (m *Manager) Install(url string) error {
	pluginPath, err := m.downloadPlugin(url, true)
	if err != nil {
		return fmt.Errorf("downloading plugin at %s: %w", url, err)
	}

	err = m.InitPlugin(pluginPath)
	if err != nil {
		return fmt.Errorf("running plugin %s after installation: %w", pluginPath, err)
	}

	return nil
}

func (m *Manager) Update(correlationID string) error {
	plgn, err := m.Plugin(correlationID)
	if err != nil {
		return fmt.Errorf("getting plugin %s: %w", plgn.info.Name, err)
	}

	if !plgn.info.Updatable {
		return fmt.Errorf("plugin %s is not updatable", plgn.info.Name)
	}

	if err != nil {
		return fmt.Errorf("while fetching store list: %w", err)
	}

	pluginInStore := store.StoreInstance.FindPluginByName(plgn.info.Name)

	url, _, _, err := pluginInStore.GetDLChecksumAndOs()
	if err != nil {
		return fmt.Errorf("while getting plugin URL of %s: %w", plgn.info.Name, err)
	}

	err = m.StopPlugin(plgn, true)
	if err != nil {
		return fmt.Errorf("stopping plugin %s: %w", plgn.info.Name, err)
	}

	// Remove the old binary file.
	pluginName := filepath.Base(plgn.Path)
	os.Remove(utils.PluginPath(plgn.Path, pluginName))

	pluginPath, err := m.downloadPlugin(url, false)
	if err != nil {
		return fmt.Errorf("downloading plugin %s: %w", plgn.info.Name, err)
	}

	err = m.InitPlugin(pluginPath)
	if err != nil {
		return fmt.Errorf("running plugin %s after update: %w", pluginPath, err)
	}

	return nil
}
