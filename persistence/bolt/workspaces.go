package bolt

import (
	"errors"
	"strings"

	"github.com/boltdb/bolt"
)

const workspacesBucket = "workspaces"
const workspacesOptionsKey = "list"
const workspacesSelectedKey = "current"

var NoSelectedWorkspaces = errors.New("there is no selected workspaces")

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
