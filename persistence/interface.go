package persistence

import (
	"errors"
	"ticktask/models"
	"ticktask/persistence/bolt"
)

var NoDataErr = errors.New("there is no data to be retrieved")

type PersistenceLayer interface {
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
}

func GetDB() PersistenceLayer {
	return bolt.GetBoltClient()
}
