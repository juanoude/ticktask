package workspace

import (
	"log"
	"slices"
	"ticktask/persistence"
	"ticktask/views"

	"github.com/spf13/cobra"
)

func init() {
	WorkspaceCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "removes a workspace",
	Long:  "Were you supposed to live with them forever?",
	Run: func(cmd *cobra.Command, args []string) {
		workspaces := persistence.GetDB().GetWorkspaces()
		selectedWorkspace := persistence.GetDB().GetSelectedWorkspace()
		selectedIndex := views.RunSelector(workspaces, "Which one you want to remove?")
		if len(workspaces) <= 1 {
			log.Fatal("you can only delete when you have more than one workspace")
		}

		persistence.GetDB().RemoveWorkspace(workspaces[selectedIndex])
		if selectedWorkspace == workspaces[selectedIndex] {
			workspaces = slices.Delete(workspaces, selectedIndex, selectedIndex+1)
			persistence.GetDB().SaveSelectedWorkspace(workspaces[0])
		}
	},
}
