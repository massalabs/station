package store

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/patrickmn/go-cache"
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
	Version string `json:"version"`
	URL     string `json:"url"`
}

type File struct {
	URL      string `json:"url"`
	Checksum string `json:"checksum"`
}

const pluginListURL = "https://raw.githubusercontent.com/massalabs/thyra-plugin-store/main/plugins.json"

const cacheExpirationMinutes = 30

const cacheDurationMinutes = 15

func FetchPluginList() ([]Plugin, error) {
	cacheDuration := cacheDurationMinutes * time.Minute // Cache the result for 15 minutes

	// Create a new cache instance with a default expiration time of 30 minutes
	c := cache.New(cacheExpirationMinutes*time.Minute, cacheDuration) //nolint:varnamelen

	// Check if the response is already cached
	if cachedResp, found := c.Get(pluginListURL); found {
		// Use the cached response if it exists
		var plugins []Plugin

		cached, ok := cachedResp.([]byte)
		if !ok {
			return nil, fmt.Errorf("casting cached JSON to bytes")
		}

		err := json.Unmarshal(cached, &plugins)
		if err != nil {
			return nil, fmt.Errorf("parsing cached JSON: %w", err)
		}

		return plugins, nil
	}

	//nolint:exhaustruct
	netClient := &http.Client{
		//nolint:gomnd
		Timeout: time.Second * 10,
	}

	// If the response is not cached, make the HTTP request
	resp, err := netClient.Get(pluginListURL) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("fetching plugin list: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body and cache the result
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	c.Set(pluginListURL, body, cacheDuration)

	// Parse the JSON response
	var plugins []Plugin

	err = json.Unmarshal(body, &plugins)
	if err != nil {
		return nil, fmt.Errorf("parsing plugin list JSON: %w", err)
	}

	return plugins, nil
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
