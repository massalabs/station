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
	m       sync.RWMutex
	status  Status
	info    *Information
}

func (p *Plugin) Information() *Information {
	p.m.RLock()
	defer p.m.RUnlock()

	return p.info
}

func (p *Plugin) SetInformation(info *Information) {
	p.m.Lock()
	defer p.m.Unlock()
	p.info = info
	p.status = Up
}

func (p *Plugin) Status() Status {
	p.m.RLock()
	defer p.m.RUnlock()

	return p.status
}

func (p *Plugin) Kill() error {
	p.m.Lock()

	p.status = ShuttingDown

	err := p.command.Process.Kill()
	if err != nil {
		p.status = Crashed
		p.m.Unlock()

		return fmt.Errorf("killing process: %w", err)
	}

	p.m.Unlock()

	err = (*p.command).Wait()
	if err.Error() != "signal: killed" {
		p.m.Lock()
		p.status = Crashed
		p.m.Unlock()

		return fmt.Errorf("killing process: unexpected wait error: got %w, want `signal: killed`", err)
	}

	return nil
}

func New(path string) (*Plugin, error) {
	p := &Plugin{status: Starting}

	p.command = exec.Command(path)

	pipe, err := p.command.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("instantiating plugin %s: initializing stdout Pipe: %w", path, err)
	}

	p.stdOut = pipe

	pipe, err = p.command.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("instantiating plugin %s: initializing stderr Pipe: %w", path, err)
	}

	p.stdErr = pipe

	err = p.command.Start()
	if err != nil {
		return nil, fmt.Errorf("instantiating plugin %s: starting command: %w", path, err)
	}

	return p, nil
}
