package workspace

import (
	"ticktask/persistence"
	"ticktask/views"

	"github.com/spf13/cobra"
)

func init() {
	WorkspaceCmd.AddCommand(selectCmd)
}

// selectCmd switches to a different workspace.
// Shows an interactive selector to choose the workspace.
// All subsequent task commands will operate on the selected workspace.
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Switch to a different workspace",
	Long:  `Opens an interactive selector to choose which workspace to use for task operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		workspaces := persistence.GetDB().GetWorkspaces()
		selectedIndex := views.RunSelector(workspaces, "Select the workspace you want to work on:")
		if selectedIndex < 0 {
			return
		}
		persistence.GetDB().SaveSelectedWorkspace(workspaces[selectedIndex])
	},
}
