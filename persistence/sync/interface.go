package sync

import (
	"log"
	"ticktask/persistence/sync/amazon"
	"ticktask/utils"
)

var bucketName = "ticktask-cli"
var fileKey = "ticktask.db"

const (
	AWSAccessKeyID     string = "aws_access_key"
	AWSSecretAccessKey string = "aws_secret_access_key"
)

type SyncClient struct {
	awsAccKeyId string
	awsSecret   string
}

type WalletDep interface {
	GetKey(key string) (string, error)
}

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

func (sc *SyncClient) Push() error {
	filePath := utils.GetInstallationPath("/data") + "/" + fileKey
	log.Println(filePath)
	config, ctx := amazon.LoadConfig(sc.awsAccKeyId, sc.awsSecret)
	var service = amazon.GetService(config, ctx)
	return service.UploadObject(bucketName, fileKey, filePath)
}

func (sc *SyncClient) Pull() error {
	filePath := utils.GetInstallationPath("/data") + "/" + fileKey
	log.Println(filePath)
	config, ctx := amazon.LoadConfig(sc.awsAccKeyId, sc.awsSecret)
	var service = amazon.GetService(config, ctx)
	return service.DownloadObject(bucketName, fileKey, filePath)
}
