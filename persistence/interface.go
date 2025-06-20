package persistence

import (
	"ticktask/models"
	"ticktask/persistence/bolt"
)

type PersistenceLayer interface {
	Get(onlyIncomplete bool) ([]models.Task, error)
	// UpdatePriority(id int, newPrio int) error
	Add(prio int, name string) error
	Complete(task models.Task) error
	Cancel(id int) error
}

func GetDB() PersistenceLayer {
	return bolt.GetBoltClient()
}
