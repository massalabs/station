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
	"github.com/massalabs/station/pkg/config"
	"github.com/massalabs/station/pkg/plugin/utils"
	"github.com/massalabs/station/pkg/store"
	"github.com/xyproto/unzip"
)

// Directory returns the plugin directory.
// Note: the plugin directory is the /plugins inside the home directory.
func Directory() string {
	homeDir, _ := config.GetConfigDir()
	pluginsDir := path.Join(homeDir, "plugins")
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
}

// NewManager instantiates a manager struct.
func NewManager() *Manager {
	//nolint:exhaust,exhaustruct
	manager := &Manager{plugins: make(map[string]*Plugin), authorNameToID: make(map[string]string)}

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

// Delete kills a plugin and remove it from the manager.
func (m *Manager) Delete(correlationID string) error {
	plgn, err := m.Plugin(correlationID)
	if err != nil {
		return fmt.Errorf("deleting plugin %s: %w", correlationID, err)
	}

	m.mutex.Lock()

	// Ignore Stop errors. We want to delete the plugin anyway
	//nolint:errcheck
	plgn.Stop()

	alias := Alias(plgn.info.Author, plgn.info.Name)

	delete(m.authorNameToID, alias)

	delete(m.plugins, correlationID)

	m.mutex.Unlock()

	err = os.RemoveAll(filepath.Dir(plgn.BinPath))
	if err != nil {
		return fmt.Errorf("deleting plugin %s: %w", correlationID, err)
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
	pluginDir := Directory()

	rootItems, err := os.ReadDir(pluginDir)
	if err != nil {
		return fmt.Errorf("reading plugins directory '%s': %w", pluginDir, err)
	}

	for _, rootItem := range rootItems {
		if rootItem.IsDir() {
			binPath := filepath.Join(pluginDir, rootItem.Name(), rootItem.Name())

			err = m.InitPlugin(binPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "WARN: while running plugin %s: %s.\n", rootItem.Name(), err)
				fmt.Fprintln(os.Stderr, "This plugin will not be executed.")
			}
		}
	}

	return nil
}

func (m *Manager) Stop() {
	config.Logger.Info("Stopping all plugins...")

	for _, plugin := range m.plugins {
		err := plugin.Stop()
		if err != nil {
			config.Logger.Warnf("Error while stopping plugin %s: %s", plugin.info.Name, err)
		}
	}
}

// DownloadPlugin downloads a plugin from a given URL.
// Pass isNew to false to update the plugin.
// Returns the plugin path.
func (m *Manager) DownloadPlugin(url string, isNew bool) (string, error) {
	pluginsDir := Directory()

	resp, err := grab.Get(pluginsDir, url)
	if err != nil {
		return "", fmt.Errorf("grabbing a plugin at %s: %w", url, err)
	}

	defer func() {
		err = os.Remove(resp.Filename)
		if err != nil {
			config.Logger.Errorf("deleting archive %s: %s", resp.Filename, err)
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

	err = plgn.Stop()
	if err != nil {
		return fmt.Errorf("stopping plugin %s: %w", plgn.info.Name, err)
	}

	_, err = m.DownloadPlugin(url, false)
	if err != nil {
		return fmt.Errorf("downloading plugin %s: %w", plgn.info.Name, err)
	}

	err = plgn.Start()
	if err != nil {
		return fmt.Errorf("starting plugin %s: %w", plgn.info.Name, err)
	}

	err = plgn.SetInformation(plgn.info.URL)
	if err != nil {
		return fmt.Errorf("setting plugin %s information: %w", plgn.info.Name, err)
	}

	return nil
}
