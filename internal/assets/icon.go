package assets

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed icon.png
var iconPNG []byte

// IconPath returns the path to a cached copy of the menu bar icon.
// The asset is a 36x36px (@2x) square template image per Apple HIG
// (18x18pt menu bar extras).
func IconPath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "tunneler")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, "icon.png")
	if err := os.WriteFile(path, iconPNG, 0o644); err != nil {
		return "", err
	}
	return path, nil
}
