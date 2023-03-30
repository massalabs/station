package plugin

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cavaliergopher/grab/v3"
	"github.com/massalabs/thyra/pkg/config"
	"github.com/xyproto/unzip"
)

// Directory returns the plugin directory.
// Note: the plugin directory is the /plugins inside the home directory.
func Directory() string {
	homeDir, _ := config.GetConfigDir()
	pluginsDir := path.Join(homeDir, "my_plugins")
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
func NewManager() (*Manager, error) {
	//nolint:exhaustruct
	manager := &Manager{plugins: make(map[string]*Plugin), authorNameToID: make(map[string]string)}

	err := manager.RunAll()
	if err != nil {
		return manager, fmt.Errorf("while running all plugin: %w", err)
	}

	return manager, nil
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

	alias := fmt.Sprintf("%s/%s", plgn.info.Author, plgn.info.Name)

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

// Install grabs a remote plugin from the given url and install it locally.
func (m *Manager) Install(url string) error {
	pluginsDir := Directory()

	resp, err := grab.Get(pluginsDir, url)
	if err != nil {
		return fmt.Errorf("grabbing a plugin at %s: %w", url, err)
	}

	defer func() {
		err = os.Remove(resp.Filename)
		if err != nil {
			log.Printf("deleting archive %s: %s", resp.Filename, err)
		}
	}()

	archiveName := filepath.Base(resp.Filename)
	pluginName := strings.Split(archiveName, ".zip")[0]
	pluginDirectory := filepath.Join(pluginsDir, pluginName)

	_, err = os.Stat(pluginDirectory)

	if os.IsNotExist(err) {
		err := os.MkdirAll(pluginDirectory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("creating %s plugin directory: creating folder %s: %w", archiveName, pluginDirectory, err)
		}
	} else {
		return fmt.Errorf("creating %s plugin directory: Plugin Already Exists", archiveName)
	}

	err = unzip.Extract(resp.Filename, pluginDirectory)
	if err != nil {
		return fmt.Errorf("extracting the plugin at %s: %w", resp.Filename, err)
	}

	err = m.InitPlugin(filepath.Join(pluginDirectory, pluginName))
	if err != nil {
		return fmt.Errorf("running plugin %s after installation: %w", pluginName, err)
	}

	return nil
}
