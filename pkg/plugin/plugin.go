package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gosimple/slug"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/plugin/utils"
	"github.com/massalabs/station/pkg/store"
	"github.com/xyproto/unzip"
)

//go:generate stringer -type=Status
type Status int64

const (
	Starting Status = iota
	//nolint:varnamelen
	Up
	Down
	ShuttingDown
	Crashed
)

//nolint:tagliatelle
type Information struct {
	Name        string   `json:"name"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	URL         *url.URL `json:"-"`
	APISpec     string   `json:"apispec"`
	Home        string   `json:"home"`
	Version     string   `json:"version"`
	Updatable   bool     `json:"-"`
}

type Plugin struct {
	command      *exec.Cmd
	mutex        sync.RWMutex
	status       Status
	info         *Information
	reverseProxy *httputil.ReverseProxy
	Path         string
	ID           string
	quitChan     chan bool
}

func (p *Plugin) Information() *Information {
	return p.info
}

func (p *Plugin) getInformation() (*Information, error) {
	manifestPath := filepath.Join(p.Path, "manifest.json")

	jsonObj, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading manifest file '%s': %w", manifestPath, err)
	}

	var info *Information

	err = json.Unmarshal(jsonObj, &info)

	if err != nil {
		return nil, fmt.Errorf("parsing manifest file '%s': %w", manifestPath, err)
	}

	logoPath := filepath.Join(p.Path, info.Logo)
	info.Logo = logoPath

	return info, nil
}

func (p *Plugin) Status() Status {
	return p.status
}

func Alias(pluginAuthor string, pluginName string) string {
	pluginAuthor = FormatTextForURL(pluginAuthor)
	pluginName = FormatTextForURL(pluginName)

	return fmt.Sprintf("%s/%s", pluginAuthor, pluginName)
}

func FormatTextForURL(text string) string {
	return slug.Make(text)
}

func (p *Plugin) Kill() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.status = ShuttingDown

	err := p.command.Process.Kill()
	<-p.quitChan

	if err != nil {
		p.status = Crashed

		return fmt.Errorf("killing process: %w", err)
	}

	// The wait function is called in the plugin start method
	p.status = Down

	return nil
}

func (p *Plugin) ReverseProxy() *httputil.ReverseProxy {
	return p.reverseProxy
}

type prefixWriter struct {
	w      io.Writer
	prefix string
}

func (pw prefixWriter) Write(buf []byte) (n int, err error) {
	data := []byte(pw.prefix + string(buf))
	n, err = pw.w.Write(data)

	if err != nil {
		return n, fmt.Errorf("writing logs with prefix: %w", err)
	}

	if n != len(data) {
		return n, io.ErrShortWrite
	}

	return len(buf), nil
}

func getPluginName(archiveName string) string {
	// if the archive name contains ".app", then it's a macOS app
	if strings.Contains(archiveName, ".app") {
		return strings.Split(archiveName, ".app")[0]
	}

	return strings.Split(archiveName, "_")[0]
}

func prepareBinary(pluginFilename, pluginPath string) error {
	pluginName := filepath.Base(pluginPath)
	binPath := filepath.Join(pluginPath, pluginFilename)

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		// If the plugin binary does not exist, it means that the plugin is a MacOS .app directory.
		// No need to rename the binary.
		// We extract the .app.zip file to the plugin directory
		if runtime.GOOS != "darwin" {
			return fmt.Errorf("plugin binary not found at %s", binPath)
		}

		appPath := filepath.Join(pluginPath, pluginName+".app")

		err = os.Mkdir(appPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("creating the plugin .app directory at %s: %w", appPath, err)
		}

		appZip := filepath.Join(pluginPath, pluginName+".app.zip")

		err = unzip.Extract(appZip, appPath)
		if err != nil {
			return fmt.Errorf("extracting the plugin at %s: %w", appZip, err)
		}

		return nil
	}

	newBinPath := utils.PluginPath(pluginPath, pluginName)

	err := os.Rename(filepath.Join(pluginPath, pluginFilename), newBinPath)
	if err != nil {
		return fmt.Errorf("renaming plugin %s: %w", pluginName, err)
	}

	return nil
}

// binPath returns the path to the plugin binary.
func (p *Plugin) binPath() string {
	pluginName := filepath.Base(p.Path)

	binPath := utils.PluginPath(p.Path, pluginName)
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		// Assuming that the plugin is a MacOS .app directory
		binPath = filepath.Join(p.Path, pluginName+".app", "Contents", "MacOS", pluginName)
	}

	return binPath
}

func (p *Plugin) Start() error {
	pluginName := filepath.Base(p.Path)

	logger.Debugf("Starting plugin '%s' with id %s\n", pluginName, p.ID)

	p.mutex.Lock()
	defer p.mutex.Unlock()

	status := p.Status()

	if status != Down && status != Starting {
		return fmt.Errorf("plugin is not ready to start")
	}

	p.command = exec.Command(p.binPath(), p.ID) // #nosec G204

	stdOutWriter := &prefixWriter{w: os.Stdout, prefix: fmt.Sprintf("[%s] - ", pluginName)}
	stdErrWriter := &prefixWriter{w: os.Stderr, prefix: fmt.Sprintf("[%s] Error: ", pluginName)}

	p.command.Stdout = stdOutWriter
	p.command.Stderr = stdErrWriter

	err := p.command.Start()
	if err != nil {
		return fmt.Errorf("start plugin %s: starting command: %w", pluginName, err)
	}

	// start a goroutine to wait on the command
	go func() {
		defer func() {
			p.quitChan <- true
		}()

		err := p.command.Wait()
		if err != nil && !(err.Error() == "signal: killed" || strings.Contains(err.Error(), "exit status")) {
			logger.Errorf("plugin '%s' exiting with error: %s\n", pluginName, err)

			p.status = Crashed

			return
		}

		logger.Debugf("plugin '%s' exiting without error.\n", pluginName)
	}()

	p.status = Up

	return nil
}

// Kills a plugin.
func (p *Plugin) Stop() error {
	logger.Debugf("Stopping plugin %s.\n", p.ID)

	status := p.Status()
	if status != Up && status != Crashed {
		return fmt.Errorf("plugin is not running")
	}

	return p.Kill()
}

func New(path, pluginID string) (*Plugin, error) {
	//nolint:exhaustruct
	plgn := &Plugin{
		status:   Starting,
		Path:     path,
		ID:       pluginID,
		quitChan: make(chan bool),
	}

	info, err := plgn.getInformation()
	if err != nil {
		return nil, fmt.Errorf("getting plugin information: %w", err)
	}

	plgn.info = info

	err = plgn.Start()
	if err != nil {
		return nil, fmt.Errorf("creating plugin: %w", err)
	}

	return plgn, nil
}

func (p *Plugin) SetInformation(parsedURL *url.URL) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	info, err := p.getInformation()
	if err != nil {
		return fmt.Errorf("error getting plugin information: %w", err)
	}

	info.URL = parsedURL

	isUpdatable, err := store.StoreInstance.CheckForPluginUpdates(info.Name, info.Version)
	if err != nil {
		logger.Errorf("error finding updates: %s", err)
	}

	info.Updatable = isUpdatable
	p.info = info

	return nil
}

func (p *Plugin) InitReverseProxy() {
	p.reverseProxy = httputil.NewSingleHostReverseProxy(p.info.URL)

	originalDirector := p.reverseProxy.Director
	p.reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}
}
