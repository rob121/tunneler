package watch

import (
	"context"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rob121/tunneler/internal/config"
	"github.com/rob121/tunneler/internal/tunnel"
)

const debounce = 300 * time.Millisecond

// Config watches tunnel and SSH config files and reloads the menu when they change.
type Config struct {
	Manager     *tunnel.Manager
	MenuChanged func()
}

// Run watches ~/.tunnels and ~/.ssh/config until ctx is cancelled.
func Run(ctx context.Context, cfg Config) {
	tunnelsPath, err := config.FilePath()
	if err != nil {
		log.Printf("watch: tunnels path: %v", err)
		return
	}

	sshConfigPath, err := config.SSHConfigPath()
	if err != nil {
		log.Printf("watch: ssh config path: %v", err)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("watch: %v", err)
		return
	}
	defer watcher.Close()

	for _, path := range []string{tunnelsPath, sshConfigPath} {
		if err := addWatch(watcher, path); err != nil {
			log.Printf("watch: %s: %v", path, err)
		}
	}

	log.Printf("watching %s and %s", tunnelsPath, sshConfigPath)

	var (
		mu     sync.Mutex
		timer  *time.Timer
		reload struct {
			tunnels   bool
			sshConfig bool
		}
	)

	schedule := func() {
		mu.Lock()
		defer mu.Unlock()

		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(debounce, func() {
			mu.Lock()
			doTunnels := reload.tunnels
			doSSH := reload.sshConfig
			reload.tunnels = false
			reload.sshConfig = false
			mu.Unlock()

			if doTunnels {
				if err := cfg.Manager.Reload(); err != nil {
					log.Printf("reload %s: %v", tunnelsPath, err)
				} else {
					log.Printf("reloaded %s", tunnelsPath)
				}
			}
			if doSSH {
				log.Printf("updated %s (applies on next tunnel start)", sshConfigPath)
			}
			if (doTunnels || doSSH) && cfg.MenuChanged != nil {
				cfg.MenuChanged()
			}
		})
	}

	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("watch: %v", err)
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if !event.Has(fsnotify.Write | fsnotify.Create | fsnotify.Rename | fsnotify.Chmod) {
				continue
			}

			switch filepath.Clean(event.Name) {
			case filepath.Clean(tunnelsPath):
				mu.Lock()
				reload.tunnels = true
				mu.Unlock()
				schedule()
			case filepath.Clean(sshConfigPath):
				mu.Lock()
				reload.sshConfig = true
				mu.Unlock()
				schedule()
			}

			// Some editors replace the file; re-add the watch if needed.
			if event.Has(fsnotify.Remove) {
				_ = addWatch(watcher, event.Name)
			}
		}
	}
}

func addWatch(w *fsnotify.Watcher, path string) error {
	if err := w.Add(path); err == nil {
		return nil
	}
	return w.Add(filepath.Dir(path))
}
