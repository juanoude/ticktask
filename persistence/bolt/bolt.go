package bolt

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"ticktask/models"
	"time"

	"github.com/boltdb/bolt"
)

type BoltClient struct{}

func GetBoltClient() *BoltClient {
	return &BoltClient{}
}

func (client *BoltClient) Open() *bolt.DB {
	// create dir if doesn't exist
	currentUser, err := user.Current()
	path := currentUser.HomeDir
	path = path + "/.ticktask/data"
	err = os.MkdirAll(path, os.ModePerm)
	db, err := bolt.Open(path+"/ticktask.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("error obtaining db lock")
	}
	return db
}

func (client *BoltClient) Get(onlyIncomplete bool) ([]models.Task, error) {
	db := client.Open()
	defer db.Close()
	var results []models.Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket := checkoutBucket(tx, "Main")
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			decodedValues := strings.Split(string(value), "::")
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

		if !onlyIncomplete {
			bucket := checkoutBucket(tx, "Main-Done")
			cursor := bucket.Cursor()
			for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
				decodedValues := strings.Split(string(value), "::")
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
			}
		}

		return nil
	})

	return results, err
}

func (client *BoltClient) Add(prio int, name string) error {
	db := client.Open()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := checkoutBucket(tx, "Main")
		id, _ := bucket.NextSequence()
		err := bucket.Put(itob(id), []byte(fmt.Sprintf("%d :: %s", prio, name)))
		return err
	})

	return err
}

func (client *BoltClient) Cancel(id int) error {
	db := client.Open()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := checkoutBucket(tx, "Main")
		err := bucket.Delete(itob(uint64(id)))
		return err
	})

	return err
}

func (client *BoltClient) Complete(task models.Task) error {
	db := client.Open()
	defer db.Close()
	// Deleting from main bucket
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := checkoutBucket(tx, "Main")
		err := bucket.Delete(itob(uint64(task.Id)))
		return err
	})

	if err != nil {
		return err
	}

	// Inserting on Done Bucket
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := checkoutBucket(tx, "Main-Done")
		id, _ := bucket.NextSequence()
		err := bucket.Put(itob(id), []byte(fmt.Sprintf("%d :: %s", task.Priority, task.Name)))
		return err
	})

	return err
}

func checkoutBucket(tx *bolt.Tx, bucketName string) *bolt.Bucket {
	tx.CreateBucketIfNotExists([]byte(bucketName))
	return tx.Bucket([]byte(bucketName))
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
