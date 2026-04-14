package utils

import (
	"fmt"
	"math/rand"
	"ticktask/models"
	"time"
)

// GetRandom selects and returns a random element from the provided slice.
// Uses the current time as the random seed for variety.
// Panics if the slice is empty.
func GetRandom(options []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(options))
	return options[randomIndex]
}

// StringifyTasks converts a slice of Task objects into human-readable strings.
// Each task is formatted as "priority -> name" (e.g., "1 -> Buy groceries").
// Used for displaying tasks in the interactive selector UI.
func StringifyTasks(tasks []models.Task) []string {
	var stringifiedTasks []string
	for _, v := range tasks {
		stringifiedTasks = append(stringifiedTasks, fmt.Sprintf("%d -> %s", v.Priority, v.Name))
	}

	return stringifiedTasks
}

// SafeArgsIndex safely retrieves an element from a slice by index.
// Returns the item and true if successful, or zero value and false if the index is out of bounds.
// This is a generic function that works with any slice type.
// Used for safely accessing command-line arguments without panicking.
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
