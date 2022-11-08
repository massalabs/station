package pluginManager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
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

var wg sync.WaitGroup

func (manager *PluginManager) List() []PluginItem {
	return manager.plugins
}

func New(hostPort int, hostTLSPort int) (*PluginManager, error) {
	log.Println("Plugin Manager initialization")

	manager := PluginManager{HostPort: hostPort, HostTLSPort: hostTLSPort}
	path, _ := os.Getwd()
	pluginDir := filepath.Join(path, "plugins")

	startPort := 4200

	//check plugin dir exists
	if _, err := os.Stat(pluginDir); err != nil {
		errMsg := "unable to find pluginDir: " + pluginDir
		log.Println(errMsg)
		return &manager, errors.New(errMsg)
	}

	// Iterate on plugins
	items, _ := ioutil.ReadDir(pluginDir)
	for _, item := range items {
		if item.IsDir() {
			pluginPath := filepath.Join(pluginDir, item.Name())
			manifest, err := DetectPlugin(pluginPath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			var plugin PluginItem
			plugin.Manifest = *manifest
			plugin.Path = pluginPath
			plugin.Port = startPort
			startPort++
			// Start plugin as go routine
			wg.Add(1)
			go StartPlugin(plugin)

			manager.plugins = append(manager.plugins, plugin)
		}
	}
	return &manager, nil
}

func DetectPlugin(path string) (*PluginManifest, error) {
	manifestFileName := "manifest.json"

	//check for manifest.json
	manifestPath := filepath.Join(path, manifestFileName)
	if _, err := os.Stat(manifestPath); err != nil {
		errMsg := "unable to find manifest.json in " + path
		return nil, errors.New(errMsg)
	}
	jsonFile, err := os.Open(manifestPath)
	if err != nil {
		fmt.Println(err)
	}
	var manifest PluginManifest
	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(jsonBytes, &manifest)

	return &manifest, nil
}

func StartPlugin(plugin PluginItem) {
	log.Println("Starting plugin '" + plugin.Manifest.Name + "' on port " + strconv.Itoa(plugin.Port))

	defer wg.Done()
	cmd := exec.Command(filepath.Join(plugin.Path, plugin.Manifest.Bin), "--port", strconv.Itoa(plugin.Port), "--path", plugin.Path)

	stdout, err := cmd.CombinedOutput()
	log.Println(string(stdout))
	log.Println(err.Error())

}
