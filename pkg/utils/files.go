package utils

import (
	"os"
	"strings"
)

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !fileInfo.IsDir()
}

// TrimBasePathFromPath trims the base path prefix from the path
func TrimBasePathFromPath(basePath string, path string) string {
	return strings.TrimPrefix(path, basePath)
}
