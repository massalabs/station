package plugin

import (
	"fmt"
	weakRand "math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

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
	plugins        map[int64]*Plugin
	authorNameToID map[string]int64
}

// NewManager instantiates a manager struct.
func NewManager() (*Manager, error) {
	weakRand.Seed(time.Now().Unix())
	//nolint:exhaustruct
	manager := &Manager{plugins: make(map[int64]*Plugin), authorNameToID: make(map[string]int64)}

	err := manager.RunAll()
	if err != nil {
		return manager, fmt.Errorf("while running all plugin: %w", err)
	}

	return manager, nil
}

// ID returns the list of all the plugin id.
func (m *Manager) ID() []int64 {
	keys := make([]int64, len(m.plugins))

	i := 0

	for k := range m.plugins {
		keys[i] = k
		i++
	}

	return keys
}

// SetAlias adds an alias to an existing plugin.
// Alias can be defined during plugin register once the name and author of the plugin can be found.
//
//nolint:varnamelen
func (m *Manager) SetAlias(name string, id int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.plugins[id] == nil {
		return fmt.Errorf("while setting alias for %s: no plugin matching the given id %d", name, id)
	}

	_, exist := m.authorNameToID[name]
	if exist {
		return fmt.Errorf("while setting alias for %s: a plugin with the same alias already exists", name)
	}

	m.authorNameToID[name] = id

	return nil
}

// PluginByAlias returns a plugin from the manager using an alias.
func (m *Manager) PluginByAlias(alias string) (*Plugin, error) {
	id, exist := m.authorNameToID[alias]
	if exist {
		p, err := m.Plugin(id)
		if err != nil {
			return nil, fmt.Errorf("getting plugin by alias %w", err)
		}

		return p, nil
	}

	return nil, fmt.Errorf("plugin not found for alias %s", alias)
}

// Plugin returns a plugin from the manager.
//
//nolint:varnamelen
func (m *Manager) Plugin(id int64) (*Plugin, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	p, ok := m.plugins[id]
	if !ok {
		return nil, fmt.Errorf("no plugin matching id %d", id)
	}

	return p, nil
}

// Delete kills a plugin and remove it from the manager.
//
//nolint:varnamelen
func (m *Manager) Delete(id int64) error {
	plgn, err := m.Plugin(id)
	if err != nil {
		return fmt.Errorf("deleting plugin %d: %w", id, err)
	}

	m.mutex.Lock()

	err = plgn.Stop()
	if err != nil {
		return err
	}

	delete(m.plugins, id)
	m.mutex.Unlock()

	err = os.RemoveAll(filepath.Dir(plgn.BinPath))
	if err != nil {
		return fmt.Errorf("deleting plugin %d: %w", id, err)
	}

	return nil
}

// generateCorrelationID generate a unique correlation id.
func (m *Manager) generateCorrelationID() int64 {
	for {
		//nolint:varnamelen
		id := int64(weakRand.Int())

		_, exist := m.plugins[id]
		if exist {
			continue
		}

		return id
	}
}

// Run starts new plugin and adds it to manager.
func (m *Manager) InitPlugin(binPath string) error {
	//nolint:varnamelen
	id := m.generateCorrelationID()

	plugin, err := New(binPath, id)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	m.plugins[id] = plugin
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

	archiveName := filepath.Base(resp.Filename)
	pluginName := strings.Split(archiveName, ".zip")[0]
	pluginDirectory := filepath.Join(pluginsDir, pluginName)

	_, err = os.Stat(pluginDirectory)

	if os.IsNotExist(err) {
		err := os.MkdirAll(pluginDirectory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("creating %s plugin directory: creating folder %s: %w", archiveName, pluginDirectory, err)
		}
	}

	err = unzip.Extract(resp.Filename, pluginDirectory)
	if err != nil {
		return fmt.Errorf("extracting the plugin at %s: %w", resp.Filename, err)
	}

	err = os.Remove(resp.Filename)
	if err != nil {
		return fmt.Errorf("deleting extracted archive %s: %w", resp.Filename, err)
	}

	err = m.InitPlugin(filepath.Join(pluginDirectory, pluginName))
	if err != nil {
		return fmt.Errorf("running plugin %s after installation: %w", pluginName, err)
	}

	return nil
}
