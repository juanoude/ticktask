// Package workspace implements workspace management commands.
// Workspaces allow organizing tasks into separate groups (e.g., "work", "personal").
// Each workspace has its own task bucket in the database.
package workspace

import (
	"fmt"

	"github.com/spf13/cobra"
)

// WorkspaceCmd is the parent command for workspace management.
// Subcommands: new, list, select, move, remove
var WorkspaceCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "Manage workspaces for your tasks",
	Long: `Workspaces allow you to organize tasks into separate groups.
Each workspace maintains its own task list and completion history.

Subcommands:
  new     - Create a new workspace
  list    - Show all workspaces
  select  - Switch to a different workspace
  move    - Move tasks between workspaces
  remove  - Delete a workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing valid subcommand like new, list or select")
	},
}
