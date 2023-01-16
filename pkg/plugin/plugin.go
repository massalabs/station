package plugin

import (
	"fmt"
	"io"
	"os/exec"
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
	Description string
	Logo        []byte
	Authority   string
	APISpec     string
}

type Plugin struct {
	command *exec.Cmd
	stdOut  io.ReadCloser
	stdErr  io.ReadCloser
	mutex   sync.RWMutex
	status  Status
	info    *Information
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

func New(path string) (*Plugin, error) {
	//nolint:exhaustruct
	plgn := &Plugin{status: Starting}

	plgn.command = exec.Command(path)

	pipe, err := plgn.command.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("instantiating plugin %s: initializing stdout Pipe: %w", path, err)
	}

	plgn.stdOut = pipe

	pipe, err = plgn.command.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("instantiating plugin %s: initializing stderr Pipe: %w", path, err)
	}

	plgn.stdErr = pipe

	err = plgn.command.Start()
	if err != nil {
		return nil, fmt.Errorf("instantiating plugin %s: starting command: %w", path, err)
	}

	return plgn, nil
}
