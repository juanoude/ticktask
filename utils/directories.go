// Package utils provides utility functions used throughout the TickTask application.
// This includes directory management, random selection, and argument handling helpers.
package utils

import (
	"errors"
	"log"
	"os"
	"os/user"
)

// GetInstallationPath constructs an absolute path within the TickTask data directory (~/.ticktask).
// It takes a relative directory path (e.g., "/data", "/music/focus") and returns the full path.
// The directory is created if it doesn't exist.
//
// Example:
//
//	GetInstallationPath("/data") → "/home/user/.ticktask/data"
//	GetInstallationPath("")      → "/home/user/.ticktask"
func GetInstallationPath(relativeDir string) string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal("error obtaining logged user")
	}
	path := currentUser.HomeDir
	path = path + "/.ticktask"
	path = path + relativeDir
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal("error creating installation directories")
	}

	return path
}

// ListFilesOnDir returns a list of filenames (not directories) in the specified path.
// Returns an error if the directory cannot be read or contains no files.
// Used primarily for listing available music files in the local music directories.
func ListFilesOnDir(path string) ([]string, error) {
	fileNamesList := []string{}
	entries, err := os.ReadDir(path)
	if err != nil {
		return fileNamesList, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fileNamesList = append(fileNamesList, entry.Name())
		}
	}

	if len(fileNamesList) == 0 {
		return fileNamesList, errors.New("no files on directory")
	}

	return fileNamesList, nil
}
