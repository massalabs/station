package update

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"github.com/hashicorp/go-version"
	"github.com/massalabs/thyra/int/systray/embedded"
	"github.com/massalabs/thyra/int/systray/utils"
	"github.com/massalabs/thyra/pkg/config"
)

const (
	// updateCheckURL is the URL to the GitHub API to check for the latest release.
	updateCheckURL = "https://api.github.com/repos/massalabs/thyra/releases/latest"

	// updateCheckInterval is the interval in seconds to check for updates.
	updateCheckInterval = 1 * time.Hour

	// updateButtonLabel is the label of the menu item that tells the user that there is a new version available.
	updateButtonLabel = "Update now"

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

// Adds a menu item to the systray menu to notify the user that there is a new version available.
//
// The menu is only added if it doesn't already exist.
func addUpdateButton(app *fyne.App, systrayMenu *fyne.Menu) {
	// We check if the menu item already exists. If it does, we don't add it again.
	for _, item := range systrayMenu.Items {
		if item.Label == updateButtonLabel {
			return
		}
	}

	updateButton := fyne.NewMenuItem(updateButtonLabel,
		func() {
			OpenUpdateDialog(app)
		},
	)

	systrayMenu.Items = append(systrayMenu.Items, updateButton)
	systrayMenu.Refresh()
}

// Checks from the latest release on GitHub if there is a newer version available.
func updateCheck(app *fyne.App, systrayMenu *fyne.Menu) {
	log.Println("Checking for updates...")

	latestVersion, err := getLatestVersion()
	if err != nil {
		log.Println("Error getting last version:", err)

		return
	}

	currentVersion, err := version.NewVersion(config.Version)
	if err != nil {
		log.Println("Error parsing current version:", err)

		return
	}

	if latestVersion.GreaterThan(currentVersion) {
		log.Println("New version available:", latestVersion)
		addUpdateButton(app, systrayMenu)

		logoNotification := fyne.NewStaticResource("logo_notification", embedded.NotificationLogo)
		utils.SetSystemTrayIcon(app, logoNotification)
	}
}

// Starts a ticker that checks for updates every updateCheckInterval.
func StartUpdateCheck(app *fyne.App, systrayMenu *fyne.Menu) {
	// We don't want to check for updates in dev mode.
	if strings.Contains(config.Version, "dev") {
		return
	}

	// We check for updates on startup.
	updateCheck(app, systrayMenu)

	go func() {
		ticker := time.NewTicker(updateCheckInterval)
		for range ticker.C {
			updateCheck(app, systrayMenu)
		}
	}()
}
