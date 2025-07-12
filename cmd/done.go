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

type appState struct {
	tasks    []string // items on the to-do list
	cursor   int      // which to-do list item our cursor is pointing at
	Selected int      // which is selected
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Completes a task",
	Long:  `Don't you want some awesome completeness madness in your goals?`,
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
		persistence.GetDB().Complete(tasks[selectedIndex], workspace)
	},
}
