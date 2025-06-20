package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version",
	Long:  `All software has versions. This is ticktack's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Tick Task Focus Nest v0.0.0")
	},
}
