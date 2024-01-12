package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/logger"
)

const (
	// updateCheckURL is the URL to the GitHub API to check for the latest release.
	updateCheckURL = "https://api.github.com/repos/massalabs/station/releases/latest"

	// httpTimeout is the timeout for the HTTP client.
	httpTimeout = 5 * time.Second
)

// We only cares about the 'tag_name' field, so we create a struct with only that field.
type release struct {
	TagName string `json:"tag_name"`
}

// Gets the latest version from GitHub.
//
// Performs a GET request to the GitHub API and decodes the response into a release struct.
func getLatestVersion() (*version.Version, error) {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, updateCheckURL, strings.NewReader(""))
	if err != nil {
		return nil, fmt.Errorf("error creating GET request for latest version: %w", err)
	}

	httpClient := &http.Client{Timeout: httpTimeout, Transport: nil, Jar: nil, CheckRedirect: nil}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error while performing GET request for latest version: %w", err)
	}

	defer resp.Body.Close()

	var release release

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	latestVersion, err := version.NewVersion(release.TagName)
	if err != nil {
		return nil, fmt.Errorf("error parsing latest version: %w", err)
	}

	return latestVersion, nil
}

// Checks from the latest release on GitHub if there is a newer version available.
func Check() bool {
	logger.Debug("Checking for updates...")

	latestVersion, err := getLatestVersion()
	if err != nil {
		logger.Errorf("Error getting last version:%s", err)

		return false
	}

	currentVersion, err := version.NewVersion(config.Version)
	if err != nil {
		logger.Errorf("Error getting current version:%s", err)

		return false
	}

	logger.Debugf("Latest version: %s; Current version: %s", latestVersion, currentVersion)

	return latestVersion.GreaterThan(currentVersion)
}
