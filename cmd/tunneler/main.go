package main

import (
	"log"

	"github.com/rob121/tunneler/internal/config"
	"github.com/rob121/tunneler/internal/logging"
	"github.com/rob121/tunneler/internal/tunnel"
	"github.com/rob121/tunneler/internal/ui"
)

func main() {
	logging.Setup()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	for i, t := range cfg.Tunnels {
		if t.Name == "" || t.SSHHost == "" || t.Dynamic <= 0 {
			log.Printf("warning: tunnel [%d] incomplete (name=%q ssh_host=%q dynamic=%d)", i+1, t.Name, t.SSHHost, t.Dynamic)
		}
	}
	if len(cfg.Tunnels) == 0 {
		log.Print("warning: no tunnels defined")
	} else if path, err := config.FilePath(); err == nil {
		log.Printf("loaded %d tunnel(s) from %s", len(cfg.Tunnels), path)
	}

	manager := tunnel.NewManager(cfg)
	ui.Run(manager)
}
