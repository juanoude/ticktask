package cmd

import (
	"log"
	"strconv"
	"ticktask/cmd/workspace"
	"ticktask/persistence"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Print the current version",
	Long:  `All software has versions. This is ticktack's`,
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
		persistence.GetDB().Add(int(priority), name, workspace)
	},
}
