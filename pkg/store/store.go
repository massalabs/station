package store

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/logger"
)

//nolint:tagliatelle
type Plugin struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Assets      struct {
		Windows    File `json:"windows"`
		Linux      File `json:"linux"`
		MacosArm64 File `json:"macos-arm64"`
		MacosAmd64 File `json:"macos-amd64"`
	} `json:"assets"`
	Version             string `json:"version"`
	URL                 string `json:"url"`
	MassaStationVersion string `json:"massaStationVersion"`
	IsCompatible        bool   `json:"-"`
}

type File struct {
	URL      string `json:"url"`
	Checksum string `json:"checksum"`
}

//nolint:gochecknoglobals,exhaustruct
var StoreInstance = &Store{}

type Store struct {
	Plugins []Plugin
	mutex   sync.RWMutex
}

const (
	pluginListURL = "https://massa-station-assets.s3.eu-west-3.amazonaws.com/plugins/plugins.json"
)

func NewStore() error {
	err := StoreInstance.FetchPluginList()
	if err != nil {
		return fmt.Errorf("while fetching plugin list: %w", err)
	}

	go StoreInstance.FetchStorePeriodically()

	return nil
}

func (s *Store) FetchPluginList() error {
	//nolint:exhaustruct
	netClient := &http.Client{
		//nolint:gomnd
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Get(pluginListURL) //nolint:noctx
	if err != nil {
		return fmt.Errorf("fetching plugin list: %w", err)
	}
	defer resp.Body.Close()

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

	for index := range plugins {
		isCompatible, err := plugins[index].IsPluginCompatible()
		if err != nil {
			return fmt.Errorf("checking if plugin is compatible: %w", err)
		}

		plugins[index].IsCompatible = isCompatible
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
			logger.Errorf("while fetching plugin list: %s", err)
		}

		logger.Debugf("Fetched plugin list. %d plugins in store.", len(s.Plugins))
	}
}

func (s *Store) CheckForPluginUpdates(name string, vers string) (bool, error) {
	pluginInStore := s.FindPluginByName(name)

	if pluginInStore == nil {
		return false, nil
	}

	pluginVersion, err := version.NewVersion(vers)
	if err != nil {
		return false, fmt.Errorf("while parsing plugin version: %w", err)
	}

	// checks if the version is greater than the current one.
	pluginInStoreVersion, err := version.NewVersion(pluginInStore.Version)
	if err != nil {
		return false, fmt.Errorf("while parsing plugin version: %w", err)
	}

	newVersionInStore := pluginInStoreVersion.GreaterThan(pluginVersion)

	isUpdatable := newVersionInStore && pluginInStore.IsCompatible

	return isUpdatable, nil
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

func (plugin *Plugin) IsPluginCompatible() (bool, error) {
	if config.Version == "dev" {
		return true, nil
	}

	config.Version = strings.TrimSuffix(config.Version, "-dev")

	massaStationVersion, err := version.NewVersion(config.Version)
	if err != nil {
		return false, fmt.Errorf("while parsing MassaStation version: %w", err)
	}

	pluginMassaStationVersionConstraint, err := version.NewConstraint(plugin.MassaStationVersion)
	if err != nil {
		return false, fmt.Errorf("while parsing MassaStation version: %w", err)
	}

	return pluginMassaStationVersionConstraint.Check(massaStationVersion), nil
}
