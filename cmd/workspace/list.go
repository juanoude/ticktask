package workspace

import (
	"fmt"
	"ticktask/persistence"

	"github.com/spf13/cobra"
)

func init() {
	WorkspaceCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all your workspaces",
	Long:  "Here you can see the entire list",
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
				fmt.Println(fmt.Sprintf("-> %s", v))
				continue
			}

			fmt.Println(v)
		}
	},
}
