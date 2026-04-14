package bolt

import (
	"errors"
	"strings"

	"github.com/boltdb/bolt"
)

// Workspace storage constants.
// Workspaces are stored in a special "workspaces" bucket with two keys:
//   - "list": All workspace names joined by "::"
//   - "current": The currently selected workspace name
const workspacesBucket = "workspaces"
const workspacesOptionsKey = "list"
const workspacesSelectedKey = "current"

// NoSelectedWorkspaces is returned when no workspace is currently selected.
var NoSelectedWorkspaces = errors.New("there is no selected workspaces")

// GetWorkspaces returns all workspace names as a string slice.
// Returns an empty slice if no workspaces exist.
func (client *BoltClient) GetWorkspaces() []string {
	db := client.Open()
	defer db.Close()
	var results []string
	db.View(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, workspacesBucket)
		if err != nil {
			return nil
		}

		value := bucket.Get([]byte(workspacesOptionsKey))
		if value == nil || len(string(value)) == 0 {
			return nil
		}
		results = strings.Split(string(value), defaultSeparator)
		return nil
	})

	return results
}

// AddWorkspace appends a new workspace name to the list.
// Workspaces are stored as a "::" separated string.
// Does not check for duplicates.
func (client *BoltClient) AddWorkspace(name string) error {
	currentList := client.GetWorkspaces()

	db := client.Open()
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, workspacesBucket)
		if err != nil {
			return err
		}

		if len(currentList) == 0 {
			return bucket.Put([]byte(workspacesOptionsKey), []byte(name))
		}

		encodedList := ""
		for i, v := range currentList {
			if i > 0 {
				encodedList += defaultSeparator + v
				continue
			}

			encodedList = v
		}

		encodedList += defaultSeparator + name
		return bucket.Put([]byte(workspacesOptionsKey), []byte(encodedList))
	})
}

// RemoveWorkspace removes a workspace by name from the list.
// Note: This only removes from the workspace list, not the task buckets.
// Task data in the workspace bucket remains in the database.
func (client *BoltClient) RemoveWorkspace(name string) error {
	currentList := client.GetWorkspaces()
	db := client.Open()
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, workspacesBucket)
		if err != nil {
			return err
		}
		if len(currentList) == 0 {
			return nil
		}

		encodedList := ""
		for _, v := range currentList {
			if v == name {
				continue
			}

			if len(encodedList) > 0 {
				encodedList += defaultSeparator + v
				continue
			}

			encodedList = v
		}

		return bucket.Put([]byte(workspacesOptionsKey), []byte(encodedList))
	})
}

// SaveSelectedWorkspace sets the currently active workspace.
// This workspace is used by default for all task operations.
func (client *BoltClient) SaveSelectedWorkspace(name string) error {
	db := client.Open()
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, workspacesBucket)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(workspacesSelectedKey), []byte(name))
	})
}

// GetSelectedWorkspace returns the currently active workspace name.
// Returns an empty string if no workspace is selected.
func (client *BoltClient) GetSelectedWorkspace() string {
	db := client.Open()
	defer db.Close()
	var current string
	db.View(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, workspacesBucket)
		if err != nil {
			return nil
		}
		currentBytes := bucket.Get([]byte(workspacesSelectedKey))
		if currentBytes != nil {
			current = string(currentBytes)
		}
		return nil
	})

	return current
}
