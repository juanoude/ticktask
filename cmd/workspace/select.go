package workspace

import (
	"ticktask/persistence"
	"ticktask/views"

	"github.com/spf13/cobra"
)

func init() {
	WorkspaceCmd.AddCommand(selectCmd)
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "selects a workspace",
	Long:  "Here you can point other commands to the respective workspace of your choice",
	Run: func(cmd *cobra.Command, args []string) {
		workspaces := persistence.GetDB().GetWorkspaces()
		selectedIndex := views.RunSelector(workspaces, "Select the workspace you want to work on:")
		persistence.GetDB().SaveSelectedWorkspace(workspaces[selectedIndex])
	},
}
