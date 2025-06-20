package cmd

import (
	"fmt"
	"log"
	"ticktask/models"
	"ticktask/persistence"
	"ticktask/views"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type appState struct {
	tasks    []string // items on the to-do list
	cursor   int      // which to-do list item our cursor is pointing at
	Selected int      // which is selected
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Completes a task",
	Long:  `Don't you want some awesome completeness madness in your goals?`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := persistence.GetDB().Get(true)
		if err != nil {
			log.Println(err.Error())
			log.Fatal("error fetching tasks")
		}

		selectedTask := views.RunSelector(tasks, "What task was masterfully done?")
		persistence.GetDB().Complete(selectedTask)
	},
}

func initialState(tasks []models.Task) appState {
	var stringifiedTasks []string
	for _, v := range tasks {
		stringifiedTasks = append(stringifiedTasks, fmt.Sprintf("%d -> %s", v.Priority, v.Name))
	}
	return appState{
		tasks:    stringifiedTasks,
		cursor:   0,
		Selected: 0,
	}
}

func (state appState) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (state appState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return state, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if state.cursor > 0 {
				state.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if state.cursor < len(state.tasks)-1 {
				state.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			state.Selected = state.cursor

			return state, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return state, nil
}

func (state appState) View() string {
	// The header
	s := "What task was completed?\n\n"

	// Iterate over our choices
	for i, choice := range state.tasks {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if state.cursor == i {
			cursor = "->" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
