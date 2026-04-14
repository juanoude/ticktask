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

// isOpen enables open-ended mode (no auto-rotation between focus/rest).
var isOpen bool

func init() {
	rootCmd.AddCommand(focusCmd)
	rootCmd.PersistentFlags().BoolVarP(&isOpen, "open", "o", false, "exceed default 25m pomodoro's time")
}

// focusCmd starts a Pomodoro-style focus timer with background music.
// First shows a task selector, then launches the timer interface.
//
// Timer modes:
//   - Focus (25 min): Green timer, focus music
//   - Rest (5 min): Red timer, rest music
//   - Generic: Blue timer, generic music
//
// Controls:
//   - Space: Toggle focus/rest
//   - Backspace: Toggle generic/chore mode
//   - q: Quit
//
// Flags:
//   -o/--open: Disable auto-rotation (open-ended session)
var focusCmd = &cobra.Command{
	Use:   "focus",
	Short: "Start a focus timer session",
	Long:  `Select a task and start a Pomodoro-style focus timer with background music.`,
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
		if selectedIndex < 0 {
			return
		}

		log.Println(fmt.Sprintf("You are focusing on: %s", tasks[selectedIndex].Name))
		views.RunCountdown(isOpen)
	},
}
