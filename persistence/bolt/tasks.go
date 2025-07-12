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

const defaultSeparator = "::"
const defaultDoneSuffix = "done"

type BoltClient struct{}

func GetBoltClient() *BoltClient {
	return &BoltClient{}
}

func (client *BoltClient) Open() *bolt.DB {
	path := utils.GetInstallationPath("/data")
	db, err := bolt.Open(path+"/ticktask.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("error obtaining db lock")
	}
	return db
}

func (client *BoltClient) Get(onlyIncomplete bool, workspace string) ([]models.Task, error) {
	db := client.Open()
	defer db.Close()
	var results []models.Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, workspace)
		log.Println(bucket)
		log.Println(err)
		if err != nil {
			return nil
		}

		cursor := bucket.Cursor()
		log.Println(cursor, "Cursor")
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

		if !onlyIncomplete {
			bucket, err := checkoutBucket(tx, fmt.Sprintf("%s-%s", workspace, defaultDoneSuffix))
			if err != nil {
				return nil
			}
			cursor := bucket.Cursor()
			recordsLimit := 5
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

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
