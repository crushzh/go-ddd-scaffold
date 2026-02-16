package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distFS embed.FS

// GetDistFS returns the embedded frontend filesystem.
// Returns an error if the dist directory does not exist (dev mode).
func GetDistFS() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}
