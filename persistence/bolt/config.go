package bolt

import (
	"errors"

	"github.com/boltdb/bolt"
)

const configBucket = "configuration"

func (client *BoltClient) StoreConfig(configKey string, configValue string) error {
	db := client.Open()
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, configBucket)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(configKey), []byte(configValue))
	})
}

func (client *BoltClient) GetConfig(configKey string) (string, error) {
	db := client.Open()
	defer db.Close()
	var result string
	finalErr := db.View(func(tx *bolt.Tx) error {
		bucket, err := checkoutBucket(tx, configBucket)
		if err != nil {
			return nil
		}

		value := bucket.Get([]byte(configKey))
		if value == nil || len(string(value)) == 0 {
			return errors.New("no value found for this config")
		}

		result = string(value)
		return nil
	})

	return result, finalErr
}
