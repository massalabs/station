package plugin

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin/utils"
	"github.com/xyproto/unzip"
)

func extractZipFilename(fileURL string) (string, error) {
	re := regexp.MustCompile(`([^/=?&]+\.zip)`)

	matches := re.FindStringSubmatch(fileURL)

	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", errors.New("zip filename not found in URL")
}

func downloadFile(url, filename string) error {
	//nolint:noctx
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	req.Header.Set("User-Agent", "MassaStation/"+config.Version)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(filename)
	if err != nil {
		//nolint:wrapcheck
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, resp.Body)

	//nolint:wrapcheck
	return err
}

// downloadPlugin downloads a plugin from a given URL.
// Pass isNew to false to update the plugin.
// Returns the plugin path.
func (m *Manager) downloadPlugin(url string, isNew bool) (string, error) {
	pluginsDir := Directory(m.configDir)

	zipFileName, err := extractZipFilename(url)
	if err != nil {
		return "", fmt.Errorf("extracting zip filename from URL %s: %w", url, err)
	}

	archivePath := filepath.Join(pluginsDir, zipFileName)

	err = downloadFile(url, archivePath)
	if err != nil {
		return "", fmt.Errorf("downloading plugin at %s: %w", url, err)
	}

	defer func() {
		err = os.Remove(archivePath)
		if err != nil {
			logger.Errorf("deleting archive %s: %s", archivePath, err)
		}
	}()

	pluginFilename := utils.PluginFileName(zipFileName)
	pluginName := getPluginName(zipFileName)
	pluginPath := filepath.Join(pluginsDir, pluginName)

	if isNew {
		_, err = os.Stat(pluginPath)
		if os.IsNotExist(err) {
			err := os.MkdirAll(pluginPath, os.ModePerm)
			if err != nil {
				return "", fmt.Errorf("creating plugin directory %s: %w", pluginPath, err)
			}
		}
	}

	err = unzip.Extract(archivePath, pluginPath)
	if err != nil {
		return "", fmt.Errorf("extracting the plugin at %s: %w", archivePath, err)
	}

	err = prepareBinary(pluginFilename, pluginPath)
	if err != nil {
		return "", fmt.Errorf("preparing plugin binary %s: %w", pluginName, err)
	}

	return pluginPath, nil
}
