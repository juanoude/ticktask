package cmd

import (
	"fmt"
	"log"
	"sort"
	"ticktask/cmd/workspace"
	"ticktask/persistence"
	"ticktask/utils"
	"ticktask/views"

	"github.com/spf13/cobra"
)

var isOpen bool

func init() {
	rootCmd.AddCommand(focusCmd)
	rootCmd.PersistentFlags().BoolVarP(&isOpen, "open", "o", false, "exceed default 25m pomodoro's time")
}

var focusCmd = &cobra.Command{
	Use:   "focus",
	Short: "Focus on a task",
	Long:  `All software has versions. This is ticktack's`,
	Run: func(cmd *cobra.Command, args []string) {
		workspace := workspace.GetSelectedWorkspace()
		tasks, err := persistence.GetDB().Get(true, workspace)
		if err != nil {
			log.Println(err)
			log.Fatal("error fetching tasks")
		}

		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Priority < tasks[j].Priority
		})

		stringifiedTasks := utils.StringifyTasks(tasks)
		selectedIndex := views.RunSelector(stringifiedTasks, "What task should be cancelled?")
		log.Println(fmt.Sprintf("You are focusing on: %s", tasks[selectedIndex].Name))
		views.RunCountdown(isOpen)
	},
}
