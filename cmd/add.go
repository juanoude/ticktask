package cmd

import (
	"fmt"
	"log"
	"strconv"
	"ticktask/cmd/workspace"
	"ticktask/persistence"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

// addCmd creates a new task in the current workspace.
// Usage: ticktask add <priority> <name>
// Priority is an integer (lower = more important).
var addCmd = &cobra.Command{
	Use:   "add [priority] [name]",
	Short: "Add a new task to the current workspace",
	Long:  `Creates a new task with the specified priority and name in the currently selected workspace.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("Priority and task name should be provided\n ticktask add [PRIORITY] [NAME]")
		}

		var priority int64
		var name string
		var error error

		priority, error = strconv.ParseInt(args[0], 10, 64)
		log.Println("Task Priority: ", priority)
		if error != nil {
			log.Fatal("priority doesn't seem to be a number")
		}

		name = args[1]
		log.Println("Task Name: ", name)
		if len(name) == 0 {
			log.Fatal("name is not provided")
		}

		workspace := workspace.GetSelectedWorkspace()
		err := persistence.GetDB().Add(int(priority), name, workspace)
		if err != nil {
			log.Fatal("error when creating task")
		}

		fmt.Println("Task created successfully")
	},
}
