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
	"strings"
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
	Logo        string
	URL         *url.URL
	APISpec     string
	Home        string
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
	exe := ""
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	//nolint:exhaustruct
	plgn := &Plugin{status: Starting, BinPath: binPath + exe, ID: pluginID, quitChan: make(chan bool)}

	err := plgn.Start()
	if err != nil {
		return nil, fmt.Errorf("creating plugin: %w", err)
	}

	return plgn, nil
}
