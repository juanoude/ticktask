package persistence

import (
	"errors"
	"ticktask/models"
	"ticktask/persistence/bolt"
	"ticktask/persistence/gkeyring"
	"ticktask/persistence/sync"
)

var NoDataErr = errors.New("there is no data to be retrieved")

const (
	AWSRegionConfig     string = "aws_region"
	AWSBucketNameConfig string = "aws_bucket_name"
)

type PersistenceLayer interface {
	//----- TODO
	Get(onlyIncomplete bool, workspace string) ([]models.Task, error)
	// UpdatePriority(id int, newPrio int) error
	Add(prio int, name string, workspace string) error
	Complete(task models.Task, workspace string) error
	Cancel(id int, workspace string) error

	//----- WORKSPACES
	GetWorkspaces() []string
	AddWorkspace(name string) error
	RemoveWorkspace(name string) error
	SaveSelectedWorkspace(name string) error
	GetSelectedWorkspace() string

	//----- CONFIGURATION
	StoreConfig(key string, value string) error
	GetConfig(key string) (string, error)
}

type WalletLayer interface {
	StoreKey(key string, value string) error
	GetKey(key string) (string, error)
}

type SyncLayer interface {
	Push() error
	Pull() error
}

func GetDB() PersistenceLayer {
	return bolt.GetBoltClient()
}

func GetWallet() WalletLayer {
	return gkeyring.GetWallet()
}

func GetSync() SyncLayer {
	wallet := GetWallet()
	return sync.GetSync(wallet)
}
