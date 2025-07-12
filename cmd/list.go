package cmd

import (
	"fmt"
	"log"
	"sort"
	"ticktask/cmd/workspace"
	"ticktask/models"
	"ticktask/persistence"

	"github.com/spf13/cobra"
)

const greenPrefix = "\033[32m"
const redPrefix = "\033[31m"
const yellowPrefix = "\033[33m"
const bluePrefix = "\033[34m"
const resetSuffix = "\033[0m"

var onlyIncomplete bool

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.PersistentFlags().BoolVarP(&onlyIncomplete, "todo", "t", false, "show only todo tasks")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `List tasks, can also display done tasks optionally`,
	Run: func(cmd *cobra.Command, args []string) {
		workspace := workspace.GetSelectedWorkspace()
		log.Println(workspace)
		tasks, err := persistence.GetDB().Get(onlyIncomplete, workspace)
		if err != nil {
			log.Fatal("error fetching tasks")
		}

		sort.Slice(tasks, func(i, j int) bool {
			return !tasks[i].IsComplete && tasks[i].Priority < tasks[j].Priority
		})
		fmt.Println()
		for _, v := range tasks {
			fmt.Printf("%s\n", renderTaskString(v))
		}
	},
}

func renderTaskString(task models.Task) string {
	if task.IsComplete {
		return fmt.Sprintf("%s[X] %d -> %s%s", redPrefix, task.Priority, task.Name, resetSuffix)
	}
	return fmt.Sprintf("%d -> %s", task.Priority, task.Name)
}
