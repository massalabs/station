package pluginmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/massalabs/thyra/pkg/config"
)

type PluginManifest struct {
	Name    string
	Bin     string
	Version string
}

type PluginItem struct {
	Manifest PluginManifest
	Path     string
	Port     int
}

type PluginManager struct {
	HostPort    int
	HostTLSPort int
	plugins     []PluginItem
}

func (manager *PluginManager) List() []PluginItem {
	return manager.plugins
}

func DetectPlugin(path string) (*PluginManifest, error) {
	manifestFileName := "manifest.json"

	manifestPath := filepath.Join(path, manifestFileName)

	if _, err := os.Stat(manifestPath); err != nil {
		return nil, errors.New("unable to find manifest.json in " + path)
	}

	jsonFile, err := os.Open(manifestPath)
	if err != nil {
		return nil, errors.New("unable to open manifest.json in " + path)
	}

	var manifest PluginManifest

	jsonBytes, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(jsonBytes, &manifest)

	if err != nil {
		return nil, errors.New("unable to decode json manifest in " + path)
	}

	return &manifest, nil
}

func StartPlugin(plugin PluginItem) {
	log.Println("Starting plugin '" + plugin.Manifest.Name + "' on port " + strconv.Itoa(plugin.Port))

	cmd := exec.Command(filepath.Join(plugin.Path, plugin.Manifest.Bin),
		"--port", strconv.Itoa(plugin.Port), "--path", plugin.Path) // #nosec G204

	stdout, err := cmd.CombinedOutput()
	log.Println(string(stdout))
	log.Println(err.Error())

	defer func() {
	}()
}

func New(hostPort int, hostTLSPort int) (*PluginManager, error) {
	log.Println("Plugin Manager initialization")

	manager := PluginManager{HostPort: hostPort, HostTLSPort: hostTLSPort, plugins: []PluginItem{}}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("reading config directory '%s': %w", configDir, err)
	}

	pluginDir := filepath.Join(configDir, "plugins")

	startPort := 4200

	if _, err := os.Stat(pluginDir); err != nil {
		errMsg := "unable to find pluginDir: " + pluginDir
		log.Println(errMsg)

		return &manager, errors.New(errMsg)
	}

	// Iterate on plugins
	items, _ := os.ReadDir(pluginDir)
	for _, item := range items {
		if item.IsDir() {
			pluginPath := filepath.Join(pluginDir, item.Name())

			manifest, err := DetectPlugin(pluginPath)
			if err != nil {
				log.Println(err)

				continue
			}

			var plugin PluginItem
			plugin.Manifest = *manifest
			plugin.Path = pluginPath
			plugin.Port = startPort
			startPort++

			// Start plugin as go routine
			go StartPlugin(plugin)

			manager.plugins = append(manager.plugins, plugin)
		}
	}

	return &manager, nil
}
