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

// removeCmd deletes a workspace from the list.
// Requires at least 2 workspaces (cannot delete the last one).
// If the deleted workspace was selected, automatically selects another.
// Note: Task data in the workspace bucket is not deleted from the database.
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Delete a workspace",
	Long:  `Opens an interactive selector to choose a workspace to delete. Cannot delete the last workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		workspaces := persistence.GetDB().GetWorkspaces()
		selectedWorkspace := persistence.GetDB().GetSelectedWorkspace()
		selectedIndex := views.RunSelector(workspaces, "Which one you want to remove?")

		if selectedIndex < 0 {
			return
		}

		if len(workspaces) <= 1 {
			log.Fatal("you can only delete when you have more than one workspace")
		}

		persistence.GetDB().RemoveWorkspace(workspaces[selectedIndex])
		// If we deleted the currently selected workspace, switch to another
		if selectedWorkspace == workspaces[selectedIndex] {
			workspaces = slices.Delete(workspaces, selectedIndex, selectedIndex+1)
			persistence.GetDB().SaveSelectedWorkspace(workspaces[0])
		}
	},
}
