package cmd

import (
	"log"
	"ticktask/cmd/workspace"
	"ticktask/persistence"
	"ticktask/utils"
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
		workspace := workspace.GetSelectedWorkspace()
		tasks, err := persistence.GetDB().Get(true, workspace)
		if err != nil {
			log.Println(err)
			log.Fatal("error fetching tasks")
		}

		stringifiedTasks := utils.StringifyTasks(tasks)
		selectedIndex := views.RunSelector(stringifiedTasks, "What task should be cancelled?")
		persistence.GetDB().Cancel(tasks[selectedIndex].Id, workspace)
	},
}
