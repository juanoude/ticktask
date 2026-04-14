// Package persistence provides the data access layer for TickTask.
// It defines interfaces for database operations, secure credential storage, and cloud sync,
// along with factory functions to obtain concrete implementations.
//
// The package uses three main abstractions:
//   - PersistenceLayer: Local database operations (tasks, workspaces, config)
//   - WalletLayer: Secure credential storage using the system keyring
//   - SyncLayer: Cloud backup/restore functionality via S3
package persistence

import (
	"errors"
	"ticktask/models"
	"ticktask/persistence/bolt"
	"ticktask/persistence/gkeyring"
	"ticktask/persistence/sync"
)

// NoDataErr is returned when a query finds no matching data.
var NoDataErr = errors.New("there is no data to be retrieved")

// Configuration keys for AWS settings stored in the local database.
const (
	AWSRegionConfig     string = "aws_region"
	AWSBucketNameConfig string = "aws_bucket_name"
)

// Configuration keys for music/Navidrome settings stored in the local database.
const (
	MusicBackendConfig          string = "music_backend"           // "local" or "navidrome"
	NavidromeBaseURLConfig      string = "navidrome_base_url"      // Server URL
	NavidromeUsernameConfig     string = "navidrome_username"      // Subsonic username
	NavidromePlaylistFocusConfig   string = "navidrome_playlist_focus"   // Focus playlist name
	NavidromePlaylistRestConfig    string = "navidrome_playlist_rest"    // Rest playlist name
	NavidromePlaylistGenericConfig string = "navidrome_playlist_generic" // Generic playlist name
)

// Wallet key for Navidrome password (stored in system keyring).
const NavidromePasswordKey string = "navidrome_password"

// PersistenceLayer defines the interface for local data storage operations.
// Implementations handle tasks, workspaces, and application configuration.
// The current implementation uses BoltDB as the underlying storage engine.
type PersistenceLayer interface {
	// Task operations

	// Get retrieves tasks from the specified workspace.
	// If onlyIncomplete is true, only pending tasks are returned.
	// Otherwise, includes up to 5 most recent completed tasks.
	Get(onlyIncomplete bool, workspace string) ([]models.Task, error)

	// Add creates a new task with the given priority and name in the workspace.
	Add(prio int, name string, workspace string) error

	// Complete marks a task as done, moving it from the active bucket to the done bucket.
	Complete(task models.Task, workspace string) error

	// Cancel permanently removes a task from the workspace.
	Cancel(id int, workspace string) error

	// Workspace operations

	// GetWorkspaces returns all workspace names.
	GetWorkspaces() []string

	// AddWorkspace creates a new workspace with the given name.
	AddWorkspace(name string) error

	// RemoveWorkspace deletes a workspace by name.
	RemoveWorkspace(name string) error

	// SaveSelectedWorkspace sets the currently active workspace.
	SaveSelectedWorkspace(name string) error

	// GetSelectedWorkspace returns the currently active workspace name.
	GetSelectedWorkspace() string

	// Configuration operations

	// StoreConfig saves a key-value configuration pair.
	StoreConfig(key string, value string) error

	// GetConfig retrieves a configuration value by key.
	GetConfig(key string) (string, error)
}

// WalletLayer defines the interface for secure credential storage.
// Implementations use the operating system's keyring/keychain service
// to store sensitive data like AWS credentials.
type WalletLayer interface {
	// StoreKey saves a secret value under the given key name.
	StoreKey(key string, value string) error

	// GetKey retrieves a secret value by key name.
	GetKey(key string) (string, error)
}

// SyncLayer defines the interface for cloud backup and restore operations.
// Used to synchronize the local database with a remote S3 bucket.
type SyncLayer interface {
	// Push uploads the local database to the remote backup location.
	Push() error

	// Pull downloads the remote backup and overwrites the local database.
	Pull() error
}

// GetDB returns the default PersistenceLayer implementation (BoltDB).
// The database file is stored at ~/.ticktask/data/ticktask.db
func GetDB() PersistenceLayer {
	return bolt.GetBoltClient()
}

// GetWallet returns the default WalletLayer implementation (system keyring).
// Uses the "ticktask-cli" service name for credential storage.
func GetWallet() WalletLayer {
	return gkeyring.GetWallet()
}

// GetSync returns a SyncLayer configured with AWS credentials from the wallet.
// Fatally exits if credentials cannot be retrieved from the keyring.
func GetSync() SyncLayer {
	wallet := GetWallet()
	return sync.GetSync(wallet)
}
