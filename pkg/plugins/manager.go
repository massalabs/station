package pluginmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/massalabs/thyra/pkg/config"
)

type PluginManifest struct {
	Name    string
	Bin     string
	Version string
}

type PluginItem struct {
	Manifest    PluginManifest
	Path        string
	Port        int
	stopChannel chan bool
}

type PluginManager struct {
	HostPort    int
	HostTLSPort int
	plugins     []PluginItem
	wg          sync.WaitGroup
}

func (manager *PluginManager) List() []PluginItem {
	return manager.plugins
}

func (manager *PluginManager) StopPlugins() {
	for _, plugin := range manager.plugins {
		log.Println("Stopping plugin " + plugin.Manifest.Name)
		plugin.stopChannel <- true
	}

	manager.wg.Wait()
}

func DetectPlugin(path string) (*PluginManifest, error) {
	manifestFileName := "manifest.json"

	manifestPath := filepath.Join(path, manifestFileName)

	if _, err := os.Stat(manifestPath); err != nil {
		return nil, errors.New("unable to find manifest.json in " + path + ": " + err.Error())
	}

	jsonFile, err := os.Open(manifestPath)
	if err != nil {
		return nil, errors.New("unable to open manifest.json in " + path + ": " + err.Error())
	}

	var manifest PluginManifest

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.New("unable to read manifest.json in " + path + ": " + err.Error())
	}

	err = json.Unmarshal(jsonBytes, &manifest)

	if err != nil {
		return nil, errors.New("unable to decode json manifest in " + path + ": " + err.Error())
	}

	return &manifest, nil
}

func StartPlugin(plugin PluginItem, waitG *sync.WaitGroup) {
	log.Println("Starting plugin '" + plugin.Manifest.Name + "' on port " + strconv.Itoa(plugin.Port))

	// Launching system command could be a security issue. This topic should be tackle, and security layers should be added
	cmd := cmd.NewCmd(filepath.Join(plugin.Path, plugin.Manifest.Bin),
		"--port", strconv.Itoa(plugin.Port), "--path", plugin.Path)
	cmd.Start() // non-blocking

	ticker := time.NewTicker(time.Second)

	// Print new lines of stdout and stdErr every seconds
	go func() {
		stdOutIdx := 0
		stdErrIdx := 0

		for range ticker.C {
			status := cmd.Status()
			stdOutLen := len(status.Stdout)
			stdErrLen := len(status.Stderr)

			if stdOutLen > stdOutIdx {
				logs := status.Stdout[stdOutIdx:]

				for _, stdOut := range logs {
					log.Println(plugin.Manifest.Name + ": " + stdOut)
				}

				stdOutIdx = stdOutLen
			}

			if stdErrLen > stdErrIdx {
				logs := status.Stderr[stdErrIdx:]

				for _, stdErr := range logs {
					log.Println(plugin.Manifest.Name + ": " + stdErr)
				}

				stdErrIdx = stdErrLen
			}
		}
	}()

	<-plugin.stopChannel

	err := cmd.Stop()
	if err != nil {
		log.Println("err " + err.Error())
	}

	cmd.Status()
	waitG.Done()
}

func New(hostPort int, hostTLSPort int) (*PluginManager, error) {
	log.Println("Plugin Manager initialization")

	manager := PluginManager{HostPort: hostPort, HostTLSPort: hostTLSPort, plugins: []PluginItem{}, wg: sync.WaitGroup{}}

	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("reading config directory '%s': %w", configDir, err)
	}

	pluginDir := filepath.Join(configDir, "plugins")

	startPort := 4200

	if _, err := os.Stat(pluginDir); err != nil {
		errMsg := "unable to find pluginDir: " + pluginDir
		log.Println(errMsg)

		return &manager, nil //nolint:nilerr
	}

	// Iterate on plugins
	items, err := os.ReadDir(pluginDir)
	if err != nil {
		return nil, fmt.Errorf("reading config directory '%s': %w", configDir, err)
	}

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

			plugin.stopChannel = make(chan bool)

			// Start plugin as go routine
			go StartPlugin(plugin, &manager.wg)

			manager.plugins = append(manager.plugins, plugin)

			manager.wg.Add(1)
		}
	}

	return &manager, nil
}
