package cmd

import (
	"log"
	"ticktask/persistence"
	"ticktask/views"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cancelCmd)
}

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a task",
	Long:  `All software has versions. This is ticktack's`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := persistence.GetDB().Get(true)
		if err != nil {
			log.Println(err)
			log.Fatal("error fetching tasks")
		}

		selectedTask := views.RunSelector(tasks, "What task should be cancelled?")
		persistence.GetDB().Cancel(selectedTask.Id)
	},
}
