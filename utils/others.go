package utils

import (
	"fmt"
	"math/rand"
	"ticktask/models"
	"time"
)

func GetRandom(options []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(options))
	return options[randomIndex]
}

func StringifyTasks(tasks []models.Task) []string {
	var stringifiedTasks []string
	for _, v := range tasks {
		stringifiedTasks = append(stringifiedTasks, fmt.Sprintf("%d -> %s", v.Priority, v.Name))
	}

	return stringifiedTasks
}

func SafeArgsIndex[T any](slice []T, index int) (item T, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			// Handle the panic gracefully
			ok = false
		}
	}()
	item = slice[index] // This might panic if the index is out of bounds
	ok = true
	return
}
