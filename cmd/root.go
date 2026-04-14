// Package cmd implements the CLI commands for TickTask using Cobra.
// The command hierarchy is:
//
//	ticktask              - Root command (shows help)
//	├── add               - Add a new task
//	├── list              - List tasks
//	├── done              - Mark task complete
//	├── cancel            - Cancel/delete task
//	├── focus             - Start focus timer
//	├── version           - Show version
//	├── workspaces        - Workspace management
//	│   ├── new           - Create workspace
//	│   ├── list          - List workspaces
//	│   ├── select        - Switch workspace
//	│   ├── move          - Move tasks between workspaces
//	│   └── remove        - Delete workspace
//	└── sync              - S3 backup/restore
//	    ├── config        - Configure AWS credentials
//	    ├── up            - Push to S3
//	    └── down          - Pull from S3
package cmd

import (
	"fmt"
	"os"
	"ticktask/cmd/sync"
	"ticktask/cmd/workspace"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Configuration flags (currently unused but reserved for future use).
var cfgFile string
var projectBase string
var userLicense string

func init() {
	// Register subcommand groups
	rootCmd.AddCommand(workspace.WorkspaceCmd)
	rootCmd.AddCommand(sync.SyncCmd)
}

// initConfig loads configuration from file (currently unused).
// Reserved for future viper-based configuration.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

// rootCmd is the base command when called without any subcommands.
// Displays a welcome message; use --help for available commands.
var rootCmd = &cobra.Command{
	Use:   "ticktask",
	Short: "Tick Task is a productivity tool to keep your focus sharp and prioritized",
	Long: `A task organizer with a focus timer built with
                love by juanoude in Go.
                Complete documentation is available at http://ticktask.dev`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("It's running dude!")
	},
}

// Execute runs the root command and handles errors.
// This is called by main.main() and is the entry point for the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
