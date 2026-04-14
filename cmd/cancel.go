package cmd

import (
	"log"
	"sort"
	"ticktask/cmd/workspace"
	"ticktask/persistence"
	"ticktask/utils"
	"ticktask/views"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cancelCmd)
}

// cancelCmd permanently removes a task from the workspace.
// Unlike "done", cancelled tasks are not preserved in the done bucket.
// Shows an interactive selector to choose which task to cancel.
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel and remove a task",
	Long:  `Opens an interactive selector to choose a task to permanently delete.`,
	Run: func(cmd *cobra.Command, args []string) {
		workspace := workspace.GetSelectedWorkspace()
		tasks, err := persistence.GetDB().Get(true, workspace)
		if err != nil {
			log.Println(err)
			log.Fatal("error fetching tasks")
		}

		stringifiedTasks := utils.StringifyTasks(tasks)
		selectedIndex := views.RunSelector(stringifiedTasks, "What task should be cancelled?")
		if selectedIndex < 0 {
			return
		}

		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Priority < tasks[j].Priority
		})

		persistence.GetDB().Cancel(tasks[selectedIndex].Id, workspace)
	},
}
