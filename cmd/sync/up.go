package sync

import (
	"ticktask/persistence"

	"github.com/spf13/cobra"
)

func init() {
	SyncCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "pushes your local as the latest db backup",
	Long:  "overwrites remote db with the one present in your local",
	Run: func(cmd *cobra.Command, args []string) {
		sync := persistence.GetSync()
		sync.Push()
	},
}
