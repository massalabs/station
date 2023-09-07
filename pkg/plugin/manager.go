package plugin

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cavaliergopher/grab/v3"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin/utils"
	"github.com/massalabs/station/pkg/store"
	"github.com/xyproto/unzip"
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

func (m *Manager) RemoveAlias(alias string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, exist := m.authorNameToID[alias]
	if !exist {
		return fmt.Errorf("while removing alias for %s: no plugin matching the given alias", alias)
	}

	delete(m.authorNameToID, alias)

	return nil
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
		return fmt.Errorf("deleting plugin %s: %w", correlationID, err)
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
// However, for the "massa-labs/massa-wallet" plugin, it keeps the wallet files.
// The wallet files are identified by the prefix "wallet_" and the suffix ".yaml".
func removePlugin(plugin *Plugin) error {
	dir := filepath.Dir(plugin.BinPath)

	alias := Alias(plugin.info.Author, plugin.info.Name)
	if alias != "massa-labs/massa-wallet" {
		err := os.RemoveAll(dir)
		if err != nil {
			return fmt.Errorf("removing plugin dir %s: %w", dir, err)
		}
	}

	// for the massa-wallet we want to keep all the wallet entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading plugin dir %s: %w", dir, err)
	}

	for _, entry := range entries {
		if !(strings.HasPrefix(entry.Name(), "wallet_") && strings.HasSuffix(entry.Name(), ".yaml")) {
			item := path.Join(dir, entry.Name())

			err = os.RemoveAll(item)
			if err != nil {
				return fmt.Errorf("removing file or directory %s: %w", item, err)
			}
		}
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
func (m *Manager) InitPlugin(binPath string) error {
	correlationID := m.generateCorrelationID()

	plugin, err := New(binPath, correlationID)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	m.plugins[correlationID] = plugin
	m.mutex.Unlock()

	return nil
}

// RunALL runs all the installed plugins.
func (m *Manager) RunAll() error {
	pluginDir := Directory(m.configDir)

	rootItems, err := os.ReadDir(pluginDir)
	if err != nil {
		return fmt.Errorf("reading plugins directory '%s': %w", pluginDir, err)
	}

	for _, rootItem := range rootItems {
		if rootItem.IsDir() {
			binPath := filepath.Join(pluginDir, rootItem.Name(), rootItem.Name())

			err = m.InitPlugin(binPath)
			if err != nil {
				logger.Warnf("While running plugin %s: %s.", rootItem.Name(), err)
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
		return fmt.Errorf("error while stopping plugin %s: %w", plugin.info.Name, err)
	}

	alias := Alias(plugin.info.Author, plugin.info.Name)

	err = m.RemoveAlias(alias)
	if err != nil {
		return fmt.Errorf("removing alias %s: %w", alias, err)
	}

	if unregister {
		err = m.UnregisterPlugin(plugin.ID)
		if err != nil {
			return fmt.Errorf("unregistering plugin %s: %w", plugin.ID, err)
		}
	}

	return nil
}

// DownloadPlugin downloads a plugin from a given URL.
// Pass isNew to false to update the plugin.
// Returns the plugin path.
func (m *Manager) DownloadPlugin(url string, isNew bool) (string, error) {
	pluginsDir := Directory(m.configDir)

	req, err := grab.NewRequest(pluginsDir, url)
	if err != nil {
		return "", fmt.Errorf("creating a new request for %s: %w", url, err)
	}

	req.HTTPRequest.Header.Set("User-Agent", fmt.Sprintf("MassaStation/%s", config.Version))

	resp := grab.DefaultClient.Do(req)
	if err := resp.Err(); err != nil {
		return "", fmt.Errorf("downloading plugin at %s: %w", url, err)
	}

	defer func() {
		err = os.Remove(resp.Filename)
		if err != nil {
			logger.Errorf("deleting archive %s: %s", resp.Filename, err)
		}
	}()

	archiveName := filepath.Base(resp.Filename)
	pluginFilename := utils.PluginFileName(archiveName)
	pluginName := strings.Split(archiveName, "_")[0]
	pluginDirectory := filepath.Join(pluginsDir, pluginName)
	pluginPath := utils.PluginPath(pluginDirectory, pluginName)

	if isNew {
		_, err = os.Stat(pluginDirectory)

		if os.IsNotExist(err) {
			err := os.MkdirAll(pluginDirectory, os.ModePerm)
			if err != nil {
				return "", fmt.Errorf("creating plugin directory %s: %w", pluginDirectory, err)
			}
		} else if _, err = os.Stat(pluginPath); err == nil {
			return "", fmt.Errorf("plugin %s already exists", pluginName)
		}
	}

	err = unzip.Extract(resp.Filename, pluginDirectory)
	if err != nil {
		return "", fmt.Errorf("extracting the plugin at %s: %w", resp.Filename, err)
	}

	err = os.Rename(filepath.Join(pluginDirectory, pluginFilename), pluginPath)
	if err != nil {
		return "", fmt.Errorf("renaming plugin %s: %w", pluginName, err)
	}

	return pluginPath, nil
}

// Install grabs a remote plugin from the given url and install it locally.
func (m *Manager) Install(url string) error {
	pluginPath, err := m.DownloadPlugin(url, true)
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

	pluginPath, err := m.DownloadPlugin(url, false)
	if err != nil {
		return fmt.Errorf("downloading plugin %s: %w", plgn.info.Name, err)
	}

	err = m.InitPlugin(pluginPath)
	if err != nil {
		return fmt.Errorf("running plugin %s after update: %w", pluginPath, err)
	}

	return nil
}
