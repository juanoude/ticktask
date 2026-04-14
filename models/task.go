// Package models defines the core domain types used throughout the TickTask application.
// These types represent the fundamental data structures for task management.
package models

// Task represents a single item in the user's task list.
// Tasks are organized by priority (lower numbers = higher priority) and can be
// marked as complete. Each task belongs to a workspace (stored separately in the
// persistence layer).
type Task struct {
	// Id is the unique identifier for the task within its workspace bucket.
	// Generated automatically by BoltDB's NextSequence().
	Id int

	// Priority determines the task's importance. Lower values indicate higher priority.
	// Used for sorting tasks in list views and selection prompts.
	Priority int

	// Name is the human-readable description of the task.
	Name string

	// IsComplete indicates whether the task has been finished.
	// Completed tasks are moved to a separate "done" bucket in the database.
	IsComplete bool
}
