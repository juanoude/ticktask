package workspace

import (
	"fmt"

	"github.com/spf13/cobra"
)

var WorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspaces for your tasks",
	Long:  "Here you can group your todos and time tracking into separate workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing valid subcommand like new, list or select")
	},
}
