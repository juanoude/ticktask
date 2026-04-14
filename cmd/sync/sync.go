// Package sync implements S3 backup and restore commands.
// Allows synchronizing the local BoltDB database with a remote S3 bucket
// for backup and multi-machine access.
package sync

import (
	"log"

	"github.com/spf13/cobra"
)

// SyncCmd is the parent command for sync operations.
// Subcommands: config, up, down
var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Backup and restore database to/from S3",
	Long: `Synchronize your local task database with an S3 bucket.

Subcommands:
  config - Set up AWS credentials and bucket information
  up     - Push local database to S3 (backup)
  down   - Pull database from S3 (restore)

Before using up/down, run 'ticktask sync config' to set up credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Missing subcommand like up, down or config")
		log.Println("Make sure you configured your local creds with the config command")
	},
}
