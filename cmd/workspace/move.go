package workspace

import (
	"log"
	"ticktask/persistence"
	"ticktask/views"

	"github.com/spf13/cobra"
)

func init() {
	WorkspaceCmd.AddCommand(moveCmd)
}

// moveCmd copies incomplete tasks from one workspace to another.
// Shows two interactive selectors: source workspace, then destination.
// Note: Tasks are copied, not moved - originals remain in the source workspace.
var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Copy tasks between workspaces",
	Long:  `Copies all incomplete tasks from a source workspace to a destination workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		list := persistence.GetDB().GetWorkspaces()
		selectedIndex := views.RunSelector(list, "What is the workspace you want to move tasks from?")
		if selectedIndex < 0 {
			return
		}

		origin := list[selectedIndex]

		// Remove source from destination options
		list = append(list[:selectedIndex], list[selectedIndex+1:]...)
		targetIndex := views.RunSelector(list, "What is the workspace you want to migrate your tasks to?")
		if targetIndex < 0 {
			return
		}

		target := list[targetIndex]

		tasks, err := persistence.GetDB().Get(true, origin)
		if err != nil {
			log.Println(err)
			log.Fatal("Error gathering tasks")
		}

		// Copy each task to the destination workspace
		for _, task := range tasks {
			persistence.GetDB().Add(task.Priority, task.Name, target)
		}
	},
}
