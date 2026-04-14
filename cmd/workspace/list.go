package workspace

import (
	"fmt"
	"ticktask/persistence"

	"github.com/spf13/cobra"
)

func init() {
	WorkspaceCmd.AddCommand(listCmd)
}

// listCmd displays all workspaces.
// The currently selected workspace is marked with "->".
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workspaces",
	Long:  `Shows all available workspaces. The current workspace is marked with an arrow.`,
	Run: func(cmd *cobra.Command, args []string) {
		list := persistence.GetDB().GetWorkspaces()
		selected := GetSelectedWorkspace()

		if len(list) == 0 {
			fmt.Println("You don't have any workspaces yet")
			fmt.Println("Run 'ticktask workspaces new <some_name>' to create a new one")
			return
		}

		fmt.Println("Here are the workspaces that you have: ")
		for _, v := range list {
			if v == selected {
				fmt.Printf("-> %s", v)
				fmt.Println()
				continue
			}

			fmt.Println(v)
		}
	},
}
