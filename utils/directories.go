package utils

import (
	"errors"
	"log"
	"os"
	"os/user"
)

// GetInstallationPath picks a relative dir path and return an equivalent homeDir directory path
// If the folder doesn't exist it will create a new one.
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
