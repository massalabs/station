package plugin

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cavaliergopher/grab/v3"
	"github.com/massalabs/thyra/pkg/config"
	"github.com/xyproto/unzip"
)

// Directory return the plugin directory.
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

type Manager struct {
	plugins map[int64]*Plugin
	m       sync.RWMutex
	id      int64
}

func NewManager() *Manager {
	//nolint:exhaustruct
	return &Manager{plugins: make(map[int64]*Plugin), id: 0}
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

// Plugin returns a plugin from the manager.
func (m *Manager) Plugin(id int64) *Plugin {
	p, ok := m.plugins[id]
	if !ok {
		return nil
	}

	return p
}

// Delete kill a plugin and remove it from the manager.
//
//nolint:varnamelen
func (m *Manager) Delete(id int64) error {
	m.m.Lock()

	plgn, ok := m.plugins[id]
	if !ok {
		m.m.Unlock()

		return errors.New("no plugin matching given id")
	}

	delete(m.plugins, id)
	m.m.Unlock()

	return plgn.Kill()
}

// Run starts new plugin and adds it to manager.
func (m *Manager) Run(file string) error {
	plugin, err := New(file)
	if err != nil {
		return err
	}

	m.m.Lock()
	m.id++
	m.plugins[m.id] = plugin
	m.m.Unlock()

	return nil
}

// plugins are expected to be located in a subdir inside default plugin directory.
func (m *Manager) RunAll() error {
	pluginDir := Directory()

	rootItems, err := os.ReadDir(pluginDir)
	if err != nil {
		return fmt.Errorf("running all plugins: reading plugins directory '%s': %w", pluginDir, err)
	}

	for _, rootItem := range rootItems {
		if !rootItem.IsDir() {
			err := m.Run(rootItem.Name())
			if err != nil {
				fmt.Fprintf(os.Stderr, "WARN: while running plugin %s: %s.\n", rootItem.Name(), err)
				fmt.Fprintln(os.Stderr, "This plugin will not be executed.")
			}
		}
	}

	return nil
}

// Install grab a remote plugin from the given url and install it locally.
func (m *Manager) Install(url string) error {
	resp, err := grab.Get(Directory(), url)
	if err != nil {
		return fmt.Errorf("grabing a plugin at %s: %w", url, err)
	}

	pluginDirectory := filepath.Dir(resp.Filename)

	err = unzip.Extract(resp.Filename, pluginDirectory)
	if err != nil {
		return fmt.Errorf("extracting the plugin at %s: %w", resp.Filename, err)
	}

	err = os.Remove(resp.Filename)
	if err != nil {
		return fmt.Errorf("deleting extracted archive %s: %w", resp.Filename, err)
	}

	fileName := filepath.Base(resp.Filename)

	splitUndescoreIndex := strings.Index(fileName, "_")
	if splitUndescoreIndex > 0 {
		fileName = fileName[:splitUndescoreIndex]
	}

	err = m.Run(filepath.Join(pluginDirectory, fileName))
	if err != nil {
		return fmt.Errorf("running after installation: %w", err)
	}

	return nil
}
