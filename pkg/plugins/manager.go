package pluginmanager

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
		errMsg := "unable to find manifest.json in " + path

		return nil, errors.New(errMsg)
	}

	jsonFile, err := os.Open(manifestPath)
	if err != nil {
		errMsg := "unable to open manifest.json in " + path

		return nil, errors.New(errMsg)
	}

	var manifest PluginManifest

	jsonBytes, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(jsonBytes, &manifest)

	if err != nil {
		errMsg := "unable to decode json manifest in " + path

		return nil, errors.New(errMsg)
	}

	return &manifest, nil
}

func StartPlugin(plugin PluginItem) {
	log.Println("Starting plugin '" + plugin.Manifest.Name + "' on port " + strconv.Itoa(plugin.Port))

	defer func() {
	}()

	cmd := exec.Command(filepath.Join(plugin.Path, plugin.Manifest.Bin),
		"--port", strconv.Itoa(plugin.Port), "--path", plugin.Path) // #nosec G204

	stdout, err := cmd.CombinedOutput()
	log.Println(string(stdout))
	log.Println(err.Error())

	defer func() {
		log.Println("Defer func!")
		log.Println(string(stdout))
		log.Println(err.Error())
	}()
}

func New(hostPort int, hostTLSPort int) (*PluginManager, error) {
	log.Println("Plugin Manager initialization")

	manager := PluginManager{HostPort: hostPort, HostTLSPort: hostTLSPort, plugins: []PluginItem{}}
	path, _ := os.Getwd()
	pluginDir := filepath.Join(path, "plugins")

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
