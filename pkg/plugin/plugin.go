package plugin

import (
	"fmt"
	"io"
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

	p.status = ShuttingDown

	err := p.command.Process.Kill()
	if err != nil {
		p.status = Crashed
		p.mutex.Unlock()

		return fmt.Errorf("killing process: %w", err)
	}

	p.mutex.Unlock()

	err = (*p.command).Wait()
	if err.Error() != "signal: killed" {
		p.mutex.Lock()
		p.status = Crashed
		p.mutex.Unlock()

		return fmt.Errorf("killing process: unexpected wait error: got %w, want `signal: killed`", err)
	}

	return nil
}

func (p *Plugin) ReverseProxy() *httputil.ReverseProxy {
	return p.reverseProxy
}

func New(binPath string, pluginID int64) (*Plugin, error) {
	//nolint:exhaustruct
	plgn := &Plugin{status: Starting}

	exe := ""
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	plgn.command = exec.Command(binPath+exe, strconv.FormatInt(pluginID, 10)) // #nosec G204

	pipe, err := plgn.command.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("instantiating plugin %s: initializing stdout Pipe: %w", binPath, err)
	}

	plgn.stdOut = pipe

	pipe, err = plgn.command.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("instantiating plugin %s: initializing stderr Pipe: %w", binPath, err)
	}

	plgn.stdErr = pipe

	err = plgn.command.Start()
	if err != nil {
		return nil, fmt.Errorf("instantiating plugin %s: starting command: %w", binPath, err)
	}

	return plgn, nil
}
