package sync

import (
	"ticktask/persistence"

	"github.com/spf13/cobra"
)

func init() {
	SyncCmd.AddCommand(downCmd)
}

// downCmd pulls the database backup from S3 and overwrites the local database.
// Warning: This replaces all local task data with the remote backup.
// Requires credentials to be configured via "ticktask sync config".
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "pulls latest db backup to your local",
	Long:  "overwrites your current db with the latest present on the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		sync := persistence.GetSync()
		sync.Pull()
	},
}
