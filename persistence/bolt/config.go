package bolt

import (
	"errors"

	"github.com/boltdb/bolt"
)

// configBucket is the BoltDB bucket name for application configuration.
// Stores non-sensitive key-value pairs like AWS region and bucket name.
const configBucket = "configuration"

// StoreConfig saves a configuration key-value pair to the database.
// Used for non-sensitive settings like AWS region and S3 bucket name.
// For sensitive data (credentials), use the WalletLayer instead.
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

// GetConfig retrieves a configuration value by key.
// Returns an error if the key doesn't exist or has an empty value.
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
