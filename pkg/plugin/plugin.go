package plugin

import (
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
	"strconv"
	"sync"
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

type Information struct {
	Name        string
	Author      string
	Description string
	Logo        []byte
	URL         *url.URL
	APISpec     string
}

type Plugin struct {
	command      *exec.Cmd
	mutex        sync.RWMutex
	status       Status
	info         *Information
	reverseProxy *httputil.ReverseProxy
	BinPath      string
	ID           int64
}

func (p *Plugin) Information() *Information {
	return p.info
}

func (p *Plugin) SetInformation(info *Information) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.info = info
	p.status = Up

	p.reverseProxy = httputil.NewSingleHostReverseProxy(p.info.URL)

	originalDirector := p.reverseProxy.Director
	p.reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}
}

func (p *Plugin) Status() Status {
	return p.status
}

func (p *Plugin) Kill() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.status = ShuttingDown

	err := p.command.Process.Kill()
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

func (p *Plugin) Start() error {
	pluginName := filepath.Base(p.BinPath)

	log.Printf("Starting plugin '%s' with id %d\n", pluginName, p.ID)

	p.mutex.Lock()

	status := p.Status()

	if status != Down && status != Starting {
		return fmt.Errorf("plugin is not ready to start")
	}

	p.command = exec.Command(p.BinPath, strconv.FormatInt(p.ID, 10)) // #nosec G204

	stdOut, err := p.command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("start plugin %s: initializing stdout Pipe: %w", pluginName, err)
	}

	stdErr, err := p.command.StderrPipe()
	if err != nil {
		return fmt.Errorf("start plugin %s: initializing stderr Pipe: %w", pluginName, err)
	}

	err = p.command.Start()
	if err != nil {
		return fmt.Errorf("start plugin %s: starting command: %w", pluginName, err)
	}

	bindPluginLogs(stdOut, stdErr)

	// start a goroutine to wait on the command
	go func() {
		err := p.command.Wait()
		if err != nil && err.Error() != "signal: killed" {
			log.Printf("plugin '%s' exiting with error: %s\n", pluginName, err)

			p.status = Crashed

			return
		}

		log.Printf("plugin '%s' exiting without error.\n", pluginName)
	}()

	p.status = Up

	p.mutex.Unlock()

	return nil
}

func bindPluginLogs(stdOut io.ReadCloser, stdErr io.ReadCloser) {
	// start two goroutine for stdout and stderr
	go func() {
		_, err := io.Copy(os.Stdout, stdOut)
		if err != nil {
			log.Printf("stdOut pipe error: %s\n", err)
		}
	}()

	go func() {
		_, err := io.Copy(os.Stderr, stdErr)
		if err != nil {
			log.Printf("stdErr pipe error: %s\n", err)
		}
	}()
}

// Kills a plugin.
func (p *Plugin) Stop() error {
	log.Printf("Stopping plugin %d.\n", p.ID)

	status := p.Status()
	if status != Up && status != Crashed {
		return fmt.Errorf("plugin is not running")
	}

	return p.Kill()
}

func New(binPath string, pluginID int64) (*Plugin, error) {
	exe := ""
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	//nolint:exhaustruct
	plgn := &Plugin{status: Starting, BinPath: binPath + exe, ID: pluginID}

	err := plgn.Start()
	if err != nil {
		return nil, fmt.Errorf("creating plugin: %w", err)
	}

	return plgn, nil
}
