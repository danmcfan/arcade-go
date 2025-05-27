package internal

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed assets
var Assets embed.FS

func AssetFiles() fs.FS {
	assetFiles, err := fs.Sub(Assets, "assets")
	if err != nil {
		log.Fatalf("Error getting asset files: %v", err)
	}
	return assetFiles
}
