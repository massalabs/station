package store

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/hashicorp/go-version"
)

//nolint:tagliatelle
type Plugin struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Assets      struct {
		Windows    File `json:"windows"`
		Linux      File `json:"linux"`
		MacosArm64 File `json:"macos-arm64"`
		MacosAmd64 File `json:"macos-amd64"`
	} `json:"assets"`
	Version             string `json:"version"`
	URL                 string `json:"url"`
	MassaStationVersion string `json:"massaStationVersion"`
}

type File struct {
	URL      string `json:"url"`
	Checksum string `json:"checksum"`
}

type Store struct {
	Plugins []Plugin
	mutex   sync.RWMutex
}

const (
	pluginListURL          = "https://raw.githubusercontent.com/massalabs/thyra-plugin-store/main/plugins.json"
	cacheExpirationMinutes = 30
)

func NewStore() (*Store, error) {
	//nolint:exhaustruct
	storeMassaStation := &Store{}

	err := storeMassaStation.FetchPluginList()
	if err != nil {
		return storeMassaStation, fmt.Errorf("while fetching plugin list: %w", err)
	}

	go storeMassaStation.FetchStorePeriodically()

	return storeMassaStation, nil
}

func (s *Store) FetchPluginList() error {
	//nolint:exhaustruct
	netClient := &http.Client{
		//nolint:gomnd
		Timeout: time.Second * 10,
	}

	// If the response is not cached, make the HTTP request
	resp, err := netClient.Get(pluginListURL) //nolint:noctx
	if err != nil {
		return fmt.Errorf("fetching plugin list: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body and cache the result
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	// Parse the JSON response
	var plugins []Plugin

	err = json.Unmarshal(body, &plugins)
	if err != nil {
		return fmt.Errorf("parsing plugin list JSON: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Plugins = plugins

	return nil
}

//nolint:varnamelen
func (plugin *Plugin) GetDLChecksumAndOs() (string, string, string, error) {
	pluginURL := ""
	os := runtime.GOOS
	checksum := ""

	switch os {
	case "linux":
		pluginURL = plugin.Assets.Linux.URL
		checksum = plugin.Assets.Linux.Checksum
	case "darwin":
		switch arch := runtime.GOARCH; arch {
		case "amd64":
			pluginURL = plugin.Assets.MacosAmd64.URL
			checksum = plugin.Assets.MacosAmd64.Checksum
		case "arm64":
			pluginURL = plugin.Assets.MacosArm64.URL
			checksum = plugin.Assets.MacosArm64.Checksum
		default:
			return pluginURL, os, checksum, fmt.Errorf("unsupported OS '%s' and arch '%s'", os, arch)
		}
	case "windows":
		pluginURL = plugin.Assets.Windows.URL
		checksum = plugin.Assets.Windows.Checksum
	default:
		return pluginURL, os, checksum, fmt.Errorf("unsupported OS '%s'", os)
	}

	return pluginURL, os, checksum, nil
}

func (s *Store) FetchStorePeriodically() {
	intervalMinutes := 10

	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		err := s.FetchPluginList()
		if err != nil {
			log.Printf("while fetching plugin list: %s", err)
		}

		log.Printf("Fetched plugin list. %d plugins in store.", len(s.Plugins))
	}
}

func (s *Store) CheckForPluginUpdates(name string, vers string) (bool, error) {
	pluginVersion, err := version.NewVersion(vers)
	if err != nil {
		return false, fmt.Errorf("while parsing plugin version: %w", err)
	}

	pluginInStore := s.FindPluginByName(name)
	if pluginInStore != nil {
		// If there is a plugin with the same name,
		// check if the version is greater than the current one.
		pluginInStoreVersion, err := version.NewVersion(pluginInStore.Version)
		if err != nil {
			return false, fmt.Errorf("while parsing plugin version: %w", err)
		}

		return pluginInStoreVersion.GreaterThan(pluginVersion), nil
	}

	return false, nil
}

func (s *Store) FindPluginByName(name string) *Plugin {
	// for each plugin in the plugins list, check if the name matches the name of the plugin
	for _, plugin := range s.Plugins {
		if plugin.Name == name {
			return &plugin
		}
	}

	return nil
}
