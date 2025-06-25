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
		timer := 25 * time.Minute
		if isOpen {
			timer = 9999 * time.Minute
		}
		views.RunCountdown(timer)
	},
}
