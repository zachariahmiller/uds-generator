// handler.go

package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:build/ui/*
var embeddedFiles embed.FS

// StaticFileHandler serves static files and handles SPA routing
func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new sub-filesystem from the embedded files, specifically pointing to build/ui
	contentStatic, err := fs.Sub(embeddedFiles, "build/ui")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Strip leading slash to avoid path issues on embedded file system
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" || path == "/" {
		path = "index.html" // Default to index.html for empty paths
	}

	// Check for static assets specifically
	if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".svg") || strings.HasSuffix(path, ".png") || path == "index.html" {
		// Use http.FileServer to serve the file from the embedded filesystem
		http.FileServer(http.FS(contentStatic)).ServeHTTP(w, r)
	} else {
		// For all other paths, serve index.html to support SPA routing
		http.ServeFile(w, r, "build/ui/index.html")
	}
}
