// Package sync implements the SyncLayer interface for cloud backup/restore.
// It uses AWS S3 to store and retrieve the BoltDB database file,
// enabling task synchronization across multiple machines.
package sync

import (
	"log"
	"ticktask/persistence/sync/amazon"
	"ticktask/utils"
)

// Default S3 bucket and object key for the database backup.
var bucketName = "ticktask-cli"
var fileKey = "ticktask.db"

// Keyring keys for AWS credentials.
// These are stored securely in the system keyring via WalletLayer.
const (
	AWSAccessKeyID     string = "aws_access_key"
	AWSSecretAccessKey string = "aws_secret_access_key"
)

// SyncClient implements the SyncLayer interface using AWS S3.
// Credentials are loaded from the system keyring at creation time.
type SyncClient struct {
	awsAccKeyId string
	awsSecret   string
}

// WalletDep defines the minimal interface needed to retrieve credentials.
// This allows for easier testing with mock implementations.
type WalletDep interface {
	GetKey(key string) (string, error)
}

// GetSync creates a new SyncClient with credentials from the wallet.
// Fatally exits if AWS credentials cannot be retrieved from the keyring.
// Users must run "ticktask sync config" first to set up credentials.
func GetSync(wallet WalletDep) *SyncClient {
	awsAccKeyId, errId := wallet.GetKey(AWSAccessKeyID)
	awsSecretId, errSecret := wallet.GetKey(AWSSecretAccessKey)
	if errSecret != nil || errId != nil {
		log.Fatal("error fetching aws credentials")
	}

	return &SyncClient{
		awsAccKeyId: awsAccKeyId,
		awsSecret:   awsSecretId,
	}
}

// Push uploads the local database file to S3.
// This overwrites any existing backup in the bucket.
// The local file path is ~/.ticktask/data/ticktask.db
func (sc *SyncClient) Push() error {
	filePath := utils.GetInstallationPath("/data") + "/" + fileKey
	log.Println(filePath)
	config, ctx := amazon.LoadConfig(sc.awsAccKeyId, sc.awsSecret)
	var service = amazon.GetService(config, ctx)
	return service.UploadObject(bucketName, fileKey, filePath)
}

// Pull downloads the database backup from S3 and overwrites the local file.
// Warning: This replaces all local task data with the remote backup.
func (sc *SyncClient) Pull() error {
	filePath := utils.GetInstallationPath("/data") + "/" + fileKey
	log.Println(filePath)
	config, ctx := amazon.LoadConfig(sc.awsAccKeyId, sc.awsSecret)
	var service = amazon.GetService(config, ctx)
	return service.DownloadObject(bucketName, fileKey, filePath)
}
