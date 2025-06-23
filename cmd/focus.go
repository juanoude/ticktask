package cmd

import (
	"fmt"
	"log"
	"sort"
	"ticktask/persistence"
	"ticktask/views"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(focusCmd)
}

var focusCmd = &cobra.Command{
	Use:   "focus",
	Short: "Focus on a task",
	Long:  `All software has versions. This is ticktack's`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := persistence.GetDB().Get(true)
		if err != nil {
			log.Println(err)
			log.Fatal("error fetching tasks")
		}

		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Priority < tasks[j].Priority
		})
		selectedTask := views.RunSelector(tasks, "What task should be cancelled?")
		fmt.Println(selectedTask.Name)
		views.RunCountdown(25 * time.Minute)
		// persistence.GetDB().Cancel(selectedTask.Id)
	},
}
