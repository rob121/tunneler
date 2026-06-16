package tunnel

import (
	"fmt"
	"os/exec"
	"strconv"
	"sync"

	"github.com/rob121/tunneler/internal/config"
)

type Manager struct {
	cfg   *config.Config
	procs map[string]*exec.Cmd
	mu    sync.Mutex
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg:   cfg,
		procs: make(map[string]*exec.Cmd),
	}
}

func (m *Manager) Tunnels() []config.Tunnel {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cfg.Tunnels
}

func (m *Manager) Reload() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.cfg = cfg
	m.mu.Unlock()
	return nil
}

func (m *Manager) Running(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	cmd, ok := m.procs[name]
	return ok && cmd.Process != nil
}

func (m *Manager) Toggle(t config.Tunnel) error {
	if m.Running(t.Name) {
		return m.Stop(t.Name)
	}
	return m.Start(t)
}

func (m *Manager) Start(t config.Tunnel) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.procs[t.Name]; ok {
		return nil
	}

	args := []string{
		"-N",
		"-D", strconv.Itoa(t.Dynamic),
		"-o", "ExitOnForwardFailure=yes",
		"-o", "ServerAliveInterval=60",
		"-o", "ServerAliveCountMax=3",
		t.SSHHost,
	}

	cmd := exec.Command("/usr/bin/ssh", args...)
	if err := cmd.Start(); err != nil {
		return err
	}

	m.procs[t.Name] = cmd

	go func(name string, c *exec.Cmd) {
		c.Wait()
		m.mu.Lock()
		delete(m.procs, name)
		m.mu.Unlock()
	}(t.Name, cmd)

	return nil
}

func (m *Manager) Stop(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cmd, ok := m.procs[name]
	if !ok {
		return nil
	}

	if cmd.Process != nil {
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("kill: %w", err)
		}
	}

	delete(m.procs, name)
	return nil
}
