package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Setup sends log output to ~/Library/Logs/Tunneler.log when stdout is not a
// terminal (for example when launched as a .app). When run from a terminal,
// logs continue to go to stdout.
func Setup() {
	if isTerminal(os.Stdout) {
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	path := filepath.Join(home, "Library", "Logs", "Tunneler.log")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}

	log.SetOutput(io.MultiWriter(f))
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
