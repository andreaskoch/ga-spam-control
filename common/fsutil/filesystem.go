// Package fsutil contains utility functions for filesystem operations.
package fsutil

import (
	"os"
	"strings"
)

// PathExists checks if the given path exists.
// Returns true if the file or directory with the given paths exists;
// ohterwise false.
func PathExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}

	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// IsDirectory checks if the given path is a directory.
// Returns true if the given path exists and is a directory;
// otherwise false.
func IsDirectory(path string) bool {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
