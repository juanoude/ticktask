package sync

import (
	"log"

	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the db up or down (uploads or downloads latest)",
	Long:  "This will make your local db copy either the latest by pushing the backup to the cloud or pulling the cloud version in order to overwrite your possibly outdated copy",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Missing subcommand like up, down or config")
		log.Println("Make sure you configured your local creds with the config command")
	},
}
