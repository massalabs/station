package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func httpNewRequest(t *testing.T, method string, url string, payload io.Reader) *http.Response {
	// Create a new HTTP request
	t.Helper()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	assert.NoError(t, err)

	// Set the Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	return resp
}

func buildAndSetCap(t *testing.T, version string, shouldSetCap bool) {
	t.Helper()
	// print current dir
	dir, err := os.Getwd()
	assert.NoError(t, err)
	log.Println("dir is ", dir)

	ldflags := "-X github.com/massalabs/thyra/pkg/config.Version=" + version
	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", "thyra-server", "../../cmd/thyra-server/main.go")
	log.Println(cmd.String())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	assert.NoError(t, err)

	if shouldSetCap {
		setcapCmd := exec.Command("sudo", "setcap", "CAP_NET_BIND_SERVICE=+eip", "thyra-server")
		setcapCmd.Stderr = os.Stderr
		setcapCmd.Stdout = os.Stdout
		err = setcapCmd.Run()

		assert.NoError(t, err)
	}
}

func getMassaStationVersion(t *testing.T) (string, string) {
	t.Helper()

	cmd := exec.Command("./thyra-server", "--version")
	cmd.Stderr = os.Stderr

	var stdout bytes.Buffer

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	assert.NoError(t, err)

	errorOutput := strings.TrimSpace(stderr.String())

	output := strings.TrimSpace(stdout.String())

	return output, errorOutput
}

func installPlugin(t *testing.T, previousVersion string) {
	t.Helper()
	//nolint: lll
	source := fmt.Sprintf("https://github.com/massalabs/thyra-node-manager-plugin/releases/download/%s/thyra-plugin-node-manager_linux-amd64.zip", previousVersion)

	installURL := fmt.Sprintf("http://my.massa:/plugin-manager?source=%s", source)
	resp := httpNewRequest(t, http.MethodPost, installURL, nil)

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, but got %d", http.StatusNoContent, resp.StatusCode)
	}

	time.Sleep(1 * time.Second)
}

func getPluginID(t *testing.T) string {
	t.Helper()

	resp := httpNewRequest(t, http.MethodGet, "http://my.massa/plugin-manager", nil)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Define a struct to represent the response item
	type Item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		// Other fields you want to extract
	}

	// Unmarshal the JSON response into a slice of items
	var items []Item
	err = json.Unmarshal(body, &items)
	assert.NoError(t, err)

	// Return nil if items is empty
	if len(items) == 0 {
		return ""
	}

	// Find the ID of the item with name "Node Manager"
	for _, item := range items {
		if item.Name == "Node Manager" {
			return item.ID
		}
	}

	return ""
}

func updatePlugin(t *testing.T, id string) int {
	t.Helper()

	url := fmt.Sprintf("http://my.massa/plugin-manager/%s/execute", id)

	// Create the request body
	payload := map[string]string{
		"command": "update",
	}
	jsonPayload, err := json.Marshal(payload)
	assert.NoError(t, err)

	resp := httpNewRequest(t, http.MethodPost, url, bytes.NewBuffer(jsonPayload))

	defer resp.Body.Close()

	return resp.StatusCode
}

func runUpdate(t *testing.T, previousVersion string) int {
	t.Helper()
	buildAndSetCap(t, previousVersion, true)

	var cmd *exec.Cmd

	port := "80"

	stop := make(chan struct{}) // Channel to communicate the stop signal

	go func() {
		cmd = exec.Command("./thyra-server", "--http-port", port)
		defer func() {
			err := cmd.Process.Kill()
			assert.NoError(t, err)
		}()

		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err := cmd.Start()
		assert.NoError(t, err)

		<-stop // Wait for the stop signal

		// Stop signal received, wait for the command to exit
		err = cmd.Wait()
		assert.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	idPlugin := getPluginID(t)
	if idPlugin == "" {
		installPlugin(t, previousVersion)
		idPlugin = getPluginID(t)
	}

	status := updatePlugin(t, idPlugin)

	err := cmd.Process.Kill()
	assert.NoError(t, err)

	return status
}

func deletePlugin(t *testing.T) {
	t.Helper()

	var cmd *exec.Cmd

	port := "80"

	stop := make(chan struct{}) // Channel to communicate the stop signal

	go func() {
		cmd = exec.Command("./thyra-server", "--http-port", port)
		defer func() {
			err := cmd.Process.Kill()
			assert.NoError(t, err)
		}()

		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err := cmd.Start()
		assert.NoError(t, err)

		<-stop // Wait for the stop signal

		// Stop signal received, wait for the command to exit
		err = cmd.Wait()
		assert.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	idPlugin := getPluginID(t)
	if idPlugin == "" {
		return
	}

	resp := httpNewRequest(t, http.MethodDelete, fmt.Sprintf("http://my.massa/plugin-manager/%s", idPlugin), nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	err := cmd.Process.Kill()
	assert.NoError(t, err)
}

func assertVersionSet(t *testing.T, version string) {
	t.Helper()
	buildAndSetCap(t, version, true)
	stdout, stderr := getMassaStationVersion(t)

	assert.Contains(t, stdout+stderr, version)
}

func TestCompileVersion(t *testing.T) {
	assertVersionSet(t, "0.0.1")
	assertVersionSet(t, "dev")
}

func TestUpdate(t *testing.T) {
	previousVersion := "v0.0.10"

	t.Run("MassaStation Version deprecated", func(t *testing.T) {
		status := runUpdate(t, previousVersion)
		assert.Equal(t, 400, status)
	})

	t.Run("MassaStation Version is dev", func(t *testing.T) {
		status := runUpdate(t, previousVersion)
		assert.Equal(t, 400, status)
	})

	t.Run("MassaStation Version is up to date", func(t *testing.T) {
		status := runUpdate(t, previousVersion)
		assert.Equal(t, 400, status) // to be changed to 204 when possible
	})
	deletePlugin(t)
	os.Remove("./thyra-server")
}
