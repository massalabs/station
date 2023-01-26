package plugin

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
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
	stdOut       io.ReadCloser
	stdErr       io.ReadCloser
	mutex        sync.RWMutex
	status       Status
	info         *Information
	reverseProxy *httputil.ReverseProxy
	BinPath      string
	ID           int64
}

func (p *Plugin) Information() *Information {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

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
	p.mutex.RLock()
	defer p.mutex.RUnlock()

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

	p.status = Down

	return nil
}

func (p *Plugin) ReverseProxy() *httputil.ReverseProxy {
	return p.reverseProxy
}

func (p *Plugin) Start() error {

	log.Printf("Starting plugin %d.\n", p.ID)

	status := p.Status()

	if status != Down && status != Starting {
		return fmt.Errorf("Plugin is not ready to start")
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.command = exec.Command(p.BinPath, strconv.FormatInt(p.ID, 10)) // #nosec G204

	pipe, err := p.command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("start plugin %s: initializing stdout Pipe: %w", p.BinPath, err)
	}

	p.stdOut = pipe

	pipe, err = p.command.StderrPipe()
	if err != nil {
		return fmt.Errorf("start plugin %s: initializing stderr Pipe: %w", p.BinPath, err)
	}

	p.stdErr = pipe

	err = p.command.Start()
	if err != nil {
		return fmt.Errorf("start plugin %s: starting command: %w", p.BinPath, err)
	}

	p.status = Up

	return nil
}

// Kills a plugin.
func (p *Plugin) Stop() error {
	log.Printf("Stopping plugin %d.\n", p.ID)

	status := p.Status()
	if status != Up {
		return fmt.Errorf("plugin is not running")
	}

	return p.Kill()
}

func New(binPath string, pluginID int64) (*Plugin, error) {
	//nolint:exhaustruct
	plgn := &Plugin{status: Starting}

	exe := ""
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	plgn.BinPath = binPath + exe
	plgn.ID = pluginID

	err := plgn.Start()
	if err != nil {
		return nil, fmt.Errorf("creating plugin: %w", err)
	}

	return plgn, nil
}
