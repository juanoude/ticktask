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

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "creates a workspace",
	Long:  "Here you can have a new one",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating new workspace")
		if v, ok := utils.SafeArgsIndex(args, 0); !ok || len(v) == 0 {
			fmt.Println("You must provide a workspace name")
			return
		}

		list := persistence.GetDB().GetWorkspaces()
		if len(list) == 0 {
			persistence.GetDB().AddWorkspace(defaultWorkspace)
		}

		persistence.GetDB().AddWorkspace(args[0])
	},
}
