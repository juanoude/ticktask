package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// versionCmd displays the current application version.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version",
	Long:  `Displays the current version of TickTask.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Tick Task Focus Nest v0.3.0")
	},
}
