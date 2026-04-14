// Package bolt implements the PersistenceLayer interface using BoltDB.
// BoltDB is an embedded key-value store that provides ACID transactions
// and stores all data in a single file (~/.ticktask/data/ticktask.db).
//
// Data organization:
//   - Each workspace has its own bucket for active tasks
//   - Completed tasks are stored in "{workspace}-done" buckets
//   - Tasks are encoded as "priority :: name" strings
//   - Task IDs are 8-byte big-endian encoded uint64 values
package bolt

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"strings"
	"ticktask/models"
	"ticktask/utils"
	"time"

	"github.com/boltdb/bolt"
)

// defaultSeparator is used to encode task fields (priority and name) into a single string.
const defaultSeparator = "::"

// defaultDoneSuffix is appended to workspace names to create the completed tasks bucket.
const defaultDoneSuffix = "done"

// BoltClient implements the PersistenceLayer interface using BoltDB.
// It manages tasks, workspaces, and configuration in a single database file.
type BoltClient struct{}

// GetBoltClient returns a new BoltClient instance.
// The client opens a new database connection for each operation.
func GetBoltClient() *BoltClient {
	return &BoltClient{}
}

// Open establishes a connection to the BoltDB database file.
// The database is stored at ~/.ticktask/data/ticktask.db with 0600 permissions.
// Uses a 1-second timeout to acquire the file lock.
// Fatally exits if the database cannot be opened.
func (client *BoltClient) Open() *bolt.DB {
	path := utils.GetInstallationPath("/data")
	db, err := bolt.Open(path+"/ticktask.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("error obtaining db lock")
	}
	return db
}

// Get retrieves tasks from the specified workspace.
// If onlyIncomplete is true, returns only active (pending) tasks.
// If onlyIncomplete is false, also includes up to 5 most recently completed tasks.
// Tasks are returned unsorted; the caller is responsible for ordering.
func (client *BoltClient) Get(onlyIncomplete bool, workspace string) ([]models.Task, error) {
	db := client.Open()
	defer db.Close()
	var results []models.Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, workspace)
		if err != nil {
			return nil
		}

		// Iterate through all active tasks in the workspace bucket
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			decodedValues := strings.Split(string(value), defaultSeparator)
			priority, err := strconv.ParseInt(strings.TrimSpace(decodedValues[0]), 10, 64)
			if err != nil {
				return err
			}

			name := strings.TrimSpace(decodedValues[1])
			results = append(results, models.Task{
				Id:         int(binary.BigEndian.Uint64(key)),
				Priority:   int(priority),
				Name:       name,
				IsComplete: false,
			})
		}

		// Optionally include recently completed tasks
		if !onlyIncomplete {
			bucket, err := checkoutBucket(tx, fmt.Sprintf("%s-%s", workspace, defaultDoneSuffix))
			if err != nil {
				return nil
			}
			cursor := bucket.Cursor()
			recordsLimit := 5
			// Iterate backwards to get most recent completions first
			for key, value := cursor.Last(); key != nil && recordsLimit > 0; key, value = cursor.Prev() {
				decodedValues := strings.Split(string(value), defaultSeparator)
				priority, err := strconv.ParseInt(strings.TrimSpace(decodedValues[0]), 10, 64)
				if err != nil {
					return err
				}

				name := strings.TrimSpace(decodedValues[1])
				results = append(results, models.Task{
					Id:         int(binary.BigEndian.Uint64(key)),
					Priority:   int(priority),
					Name:       name,
					IsComplete: true,
				})
				recordsLimit -= 1
			}
		}

		return nil
	})

	return results, err
}

// Add creates a new task in the specified workspace.
// The task is stored with an auto-incremented ID and encoded as "priority :: name".
func (client *BoltClient) Add(prio int, name string, workspace string) error {
	db := client.Open()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, workspace)
		if err != nil {
			return err
		}

		id, _ := bucket.NextSequence()
		err = bucket.Put(itob(id), []byte(fmt.Sprintf("%d :: %s", prio, name)))
		return err
	})

	return err
}

// Cancel permanently removes a task from the workspace by its ID.
// The task is deleted entirely, not moved to the done bucket.
func (client *BoltClient) Cancel(id int, workspace string) error {
	db := client.Open()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, workspace)
		if err != nil {
			return err
		}

		err = bucket.Delete(itob(uint64(id)))
		return err
	})

	return err
}

// Complete marks a task as finished by moving it from the active workspace bucket
// to the "{workspace}-done" bucket. The task receives a new ID in the done bucket.
func (client *BoltClient) Complete(task models.Task, workspace string) error {
	db := client.Open()
	defer db.Close()
	// Deleting from main bucket
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, workspace)
		if err != nil {
			return err
		}

		err = bucket.Delete(itob(uint64(task.Id)))
		return err
	})

	if err != nil {
		return err
	}

	// Inserting on Done Bucket
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, fmt.Sprintf("%s-%s", workspace, defaultDoneSuffix))
		if err != nil {
			return err
		}

		id, _ := bucket.NextSequence()
		err = bucket.Put(itob(id), []byte(fmt.Sprintf("%d :: %s", task.Priority, task.Name)))
		return err
	})

	return err
}

// checkoutBucket retrieves or creates a bucket with the given name.
// If the bucket doesn't exist and the transaction is writable, it creates the bucket.
// This is a helper function used by all database operations.
func checkoutBucket(tx *bolt.Tx, bucketName string) (*bolt.Bucket, error) {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket != nil {
		return bucket, nil
	}

	_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return nil, err
	}
	return tx.Bucket([]byte(bucketName)), nil
}

// itob converts a uint64 to an 8-byte big-endian byte slice.
// Used to encode task IDs as BoltDB keys, which are sorted lexicographically.
// Big-endian encoding ensures numeric ordering matches lexicographic ordering.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
