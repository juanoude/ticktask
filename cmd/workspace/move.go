package workspace

import (
	"fmt"
	"log"
	"ticktask/persistence"
	"ticktask/views"

	"github.com/spf13/cobra"
)

func init() {
	WorkspaceCmd.AddCommand(moveCmd)
}

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "migrates your todos",
	Long:  "Here your can change tasks from one workspace to another",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Here is the workspaces that you have: ")
		list := persistence.GetDB().GetWorkspaces()
		selectedIndex := views.RunSelector(list, "What is the workspace you want to move tasks from?")
		origin := list[selectedIndex]

		list = append(list[:selectedIndex], list[selectedIndex+1:]...)
		targetIndex := views.RunSelector(list, "What is the workspace you want to migrate your tasks to?")
		target := list[targetIndex]

		tasks, err := persistence.GetDB().Get(true, origin)
		if err != nil {
			log.Println(err)
			log.Fatal("Error gathering tasks")
		}

		for _, task := range tasks {
			persistence.GetDB().Add(task.Priority, task.Name, target)
		}
	},
}
