package workspace

import (
	"fmt"
	"ticktask/persistence"
	"ticktask/utils"

	"github.com/spf13/cobra"
)

func init() {
	WorkspaceCmd.AddCommand(newCmd)
}

// newCmd creates a new workspace.
// If this is the first workspace, also creates a "default" workspace.
// Usage: ticktask workspaces new <name>
var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new workspace",
	Long:  `Creates a new workspace with the given name for organizing tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		if v, ok := utils.SafeArgsIndex(args, 0); !ok || len(v) == 0 {
			fmt.Println("You must provide a workspace name")
			return
		}

		// Ensure "default" workspace exists before creating others
		list := persistence.GetDB().GetWorkspaces()
		if len(list) == 0 {
			persistence.GetDB().AddWorkspace(defaultWorkspace)
		}

		persistence.GetDB().AddWorkspace(args[0])
	},
}
