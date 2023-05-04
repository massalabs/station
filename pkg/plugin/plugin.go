package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	"github.com/hashicorp/go-version"

	"github.com/massalabs/thyra/pkg/store"
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
	BinPath      string
	ID           string
	quitChan     chan bool
}

func (p *Plugin) Information() *Information {
	return p.info
}

func (p *Plugin) getInformation() (*Information, error) {
	manifestPath := filepath.Join(filepath.Dir(p.BinPath), "manifest.json")

	jsonObj, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading manifest file '%s': %w", manifestPath, err)
	}

	var manifest *Information

	err = json.Unmarshal(jsonObj, &manifest)

	if err != nil {
		return nil, fmt.Errorf("parsing manifest file '%s': %w", manifestPath, err)
	}

	return manifest, nil
}

func (p *Plugin) SetInformation(parsedURL *url.URL) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	info, err := p.getInformation()
	if err != nil {
		return fmt.Errorf("error getting plugin information: %w", err)
	}

	info.URL = parsedURL

	isUpdatable, err := p.checkForPluginUpdates()
	if err != nil {
		return fmt.Errorf("error finding updates: %w", err)
	}

	info.Updatable = isUpdatable
	p.info = info
	p.status = Up

	p.reverseProxy = httputil.NewSingleHostReverseProxy(p.info.URL)

	originalDirector := p.reverseProxy.Director
	p.reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	return nil
}

func findPluginByName(name string, plugins []store.Plugin) *store.Plugin {
	// for each plugin in the plugins list, check if the name matches the name of the plugin
	for _, plugin := range plugins {
		if plugin.Name == name {
			return &plugin
		}
	}

	return nil
}

func (p *Plugin) checkForPluginUpdates() (bool, error) {
	// check if there is an element pluginInStore in storeitems that has the same name as rootItem.Name()
	pluginList, err := store.FetchPluginList()
	if err != nil {
		return false, fmt.Errorf("while fetching store list: %w", err)
	}

	pluginVersion, err := version.NewVersion(p.info.Version)
	if err != nil {
		return false, fmt.Errorf("while parsing plugin version: %w", err)
	}

	pluginInStore := findPluginByName(p.info.Name, pluginList)
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

func (p *Plugin) Start() error {
	pluginName := filepath.Base(p.BinPath)

	log.Printf("Starting plugin '%s' with id %s\n", pluginName, p.ID)

	p.mutex.Lock()
	defer p.mutex.Unlock()

	status := p.Status()

	if status != Down && status != Starting {
		return fmt.Errorf("plugin is not ready to start")
	}

	p.command = exec.Command(p.BinPath, p.ID) // #nosec G204

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
			log.Printf("plugin '%s' exiting with error: %s\n", pluginName, err)

			p.status = Crashed

			return
		}

		log.Printf("plugin '%s' exiting without error.\n", pluginName)
	}()

	p.status = Up

	return nil
}

// Kills a plugin.
func (p *Plugin) Stop() error {
	log.Printf("Stopping plugin %s.\n", p.ID)

	status := p.Status()
	if status != Up && status != Crashed {
		return fmt.Errorf("plugin is not running")
	}

	return p.Kill()
}

func New(binPath string, pluginID string) (*Plugin, error) {
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	//nolint:exhaustruct
	plgn := &Plugin{
		status:   Starting,
		BinPath:  binPath + ext,
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
