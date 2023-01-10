package pluginmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/go-cmd/cmd"
	"github.com/massalabs/thyra/pkg/config"
)

const (
	fileMode = 0o755
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

func (manager *PluginManager) StartPlugins() {
	for _, plugin := range manager.plugins {
		go StartPlugin(plugin, &manager.wg)

		manager.wg.Add(1)
	}
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

	cmdOptions := cmd.Options{ //nolint:exhaustruct
		Buffered:  false,
		Streaming: true,
	}
	// Launching system command could be a security issue. This topic should be tackle, and security layers should be added
	cmd := cmd.NewCmdOptions(cmdOptions, filepath.Join(plugin.Path, plugin.Manifest.Bin),
		"--port", strconv.Itoa(plugin.Port), "--path", plugin.Path)

	// Print STDOUT and STDERR lines streaming from Cmd
	printChannel := make(chan struct{})

	go func() {
		defer close(printChannel)
		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for cmd.Stdout != nil || cmd.Stderr != nil {
			select {
			case line, open := <-cmd.Stdout:
				if !open {
					cmd.Stdout = nil

					continue
				}

				log.Println(plugin.Manifest.Name + ": " + line)
			case line, open := <-cmd.Stderr:
				if !open {
					cmd.Stderr = nil

					continue
				}

				fmt.Fprintln(os.Stderr, plugin.Manifest.Name+": "+line)
			}
		}
	}()

	cmd.Start() // non-blocking

	<-plugin.stopChannel

	err := cmd.Stop()
	if err != nil {
		log.Println("err " + err.Error())
	}

	// Wait for goroutine to print everything
	<-printChannel

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
			manager.plugins = append(manager.plugins, plugin)
		}
	}

	err = manager.InstallNodeManager()
	if err != nil {
		return nil, fmt.Errorf("installing NodeManager plugin: %w", err)
	}

	return &manager, nil
}

// TEMPORARY. This code should be removed once the NodeManager plugin tests are done.
func (manager *PluginManager) InstallNodeManager() error {
	for _, plugin := range manager.plugins {
		if plugin.Manifest.Name == "Node Manager plugin" {
			return nil
		}
	}

	log.Println("NodeManager plugin not found, installing it...")

	configDir, err := config.GetConfigDir()
	if err != nil {
		return fmt.Errorf("reading config directory '%s': %w", configDir, err)
	}

	pluginDir := filepath.Join(configDir, "plugins")

	err = os.MkdirAll(pluginDir, fileMode)
	if err != nil {
		return fmt.Errorf("creating plugin directory '%s': %w", pluginDir, err)
	}

	zipPath := filepath.Join(pluginDir, "plugin.zip")

	err = downloadFile(zipPath)
	if err != nil {
		return fmt.Errorf("downloading plugin: %w", err)
	}

	err = unzip(zipPath, pluginDir)
	if err != nil {
		return err
	}

	err = os.Remove(zipPath)
	if err != nil {
		return fmt.Errorf("removing zip file '%s': %w", zipPath, err)
	}

	pluginPath := filepath.Join(pluginDir, "thyra-node-manager-plugin")

	manifest, err := DetectPlugin(pluginPath)
	if err != nil {
		log.Println(err)
	}

	err = os.Chmod(filepath.Join(pluginPath, manifest.Bin), fileMode)
	if err != nil {
		return fmt.Errorf("chmod on plugin executable '%s': %w", filepath.Join(pluginPath, manifest.Bin), err)
	}

	var plugin PluginItem
	plugin.Manifest = *manifest
	plugin.Path = pluginPath
	plugin.Port = 4222

	plugin.stopChannel = make(chan bool)
	manager.plugins = append(manager.plugins, plugin)

	return nil
}

func unzip(zipPath string, pluginDir string) error {
	cmd := exec.Command("tar", "-xf", zipPath, "-C", pluginDir)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("unzipping plugin '%s': %w", zipPath, err)
	}

	return nil
}

func downloadFile(filepath string) error {
	pluginURL, err := getDownloadURL()
	if err != nil {
		return fmt.Errorf("getting download url: %w", err)
	}

	resp, err := http.Get(pluginURL) //nolint:gosec,noctx
	if err != nil {
		return fmt.Errorf("downloading file '%s': %w", pluginURL, err)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating file '%s': %w", filepath, err)
	}

	defer func() {
		err := out.Close()
		if err != nil {
			log.Println("error closing file: " + err.Error())
		}

		err = resp.Body.Close()
		if err != nil {
			log.Println("error closing response body: " + err.Error())
		}
	}()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("writing file '%s': %w", filepath, err)
	}

	return nil
}

func getDownloadURL() (string, error) {
	pluginURL := ""

	switch os := runtime.GOOS; os { //nolint:varnamelen
	case "linux":
		pluginURL = "https://github.com/massalabs/thyra-node-manager-plugin/releases/latest/download/node-manager-plugin-linux_amd64.zip" //nolint:lll
	case "darwin":
		switch arch := runtime.GOARCH; arch {
		case "amd64":
			pluginURL = "https://github.com/massalabs/thyra-node-manager-plugin/releases/latest/download/node-manager-plugin-darwin_amd64.zip" //nolint:lll
		case "arm64":
			pluginURL = "https://github.com/massalabs/thyra-node-manager-plugin/releases/latest/download/node-manager-plugin-darwin_arm64.zip" //nolint:lll
		default:
			return "", fmt.Errorf("unsupported OS '%s' and arch '%s'", os, arch)
		}
	case "windows":
		pluginURL = "https://github.com/massalabs/thyra-node-manager-plugin/releases/latest/download/node-manager-plugin-windows_amd64.zip" //nolint:lll
	default:
		return "", fmt.Errorf("unsupported OS '%s'", os)
	}

	return pluginURL, nil
}
