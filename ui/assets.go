//go:build !dev
// +build !dev

package ui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed all:build
var _assets embed.FS

var Assets http.FileSystem

func init() {
	sub, err := fs.Sub(_assets, "build")
	if err != nil {
		log.Fatal("Unable to cd to ui/build directory:", err)
	}
	Assets = http.FS(sub)
}

func sanitizeStaticOptions() error {
	return nil
}
