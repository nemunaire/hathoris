//go:build dev
// +build dev

package ui

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"
)

var (
	Assets    http.FileSystem
	StaticDir string = "ui/"
)

func init() {
	flag.StringVar(&StaticDir, "static", StaticDir, "Directory containing static files")
}

func sanitizeStaticOptions() error {
	StaticDir, _ = filepath.Abs(StaticDir)
	if _, err := os.Stat(StaticDir); os.IsNotExist(err) {
		StaticDir, _ = filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "ui"))
		if _, err := os.Stat(StaticDir); os.IsNotExist(err) {
			return err
		}
	}
	Assets = http.Dir(StaticDir)
	return nil
}
