package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetWithExt(filename string, mediaType string) string {
	ext := mediaTypeToExt(mediaType)
	return fmt.Sprintf("%s%s", filename, ext)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func mediaTypeToExt(mediaType string) string {
	var extension string

	switch mediaType {
	case "image/jpeg":
		extension = ".jpg"
	case "image/png":
		extension = ".png"
	case "video/mp4":
		extension = ".mp4"
	default:
		extension = ".bin"
	}

	return extension
}

func generateRandomByte() ([]byte, error) {
	randomByte := make([]byte, 32)
	_, err := rand.Read(randomByte)
	if err != nil {
		return nil, err
	}

	return randomByte, nil
}
