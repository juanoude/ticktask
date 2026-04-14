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

// appState is unused but kept for potential future use.
type appState struct {
	tasks    []string
	cursor   int
	Selected int
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

// doneCmd marks a task as completed.
// Shows an interactive selector to choose which task to complete.
// The task is moved from the active bucket to the "done" bucket.
var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Mark a task as completed",
	Long:  `Opens an interactive selector to choose a task to mark as done.`,
	Run: func(cmd *cobra.Command, args []string) {
		workspace := workspace.GetSelectedWorkspace()
		tasks, err := persistence.GetDB().Get(true, workspace)
		if err != nil {
			log.Println(err.Error())
			log.Fatal("error fetching tasks")
		}

		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Priority < tasks[j].Priority
		})

		stringifiedTasks := utils.StringifyTasks(tasks)
		selectedIndex := views.RunSelector(stringifiedTasks, "What task was masterfully done?")
		if selectedIndex < 0 {
			return
		}

		persistence.GetDB().Complete(tasks[selectedIndex], workspace)
	},
}
