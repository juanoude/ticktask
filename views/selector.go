package views

import (
	"fmt"
	"os"
	"ticktask/models"

	tea "github.com/charmbracelet/bubbletea"
)

type selectorState struct {
	question string
	tasks    []string // items on the to-do list
	cursor   int      // which to-do list item our cursor is pointing at
	Selected int      // which is selected
}

func RunSelector(tasks []models.Task, question string) models.Task {
	p := tea.NewProgram(initSelector(tasks, question))
	result, err := p.Run()
	if err != nil {
		fmt.Printf("Oops, there's been an error: %v", err)
		os.Exit(1)
	}

	finalResult := result.(selectorState)
	return tasks[finalResult.Selected]
}

func initSelector(tasks []models.Task, questionString string) selectorState {
	var stringifiedTasks []string
	for _, v := range tasks {
		stringifiedTasks = append(stringifiedTasks, fmt.Sprintf("%d -> %s", v.Priority, v.Name))
	}
	return selectorState{
		question: questionString,
		tasks:    stringifiedTasks,
		cursor:   0,
		Selected: 0,
	}
}

func (state selectorState) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (state selectorState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (state selectorState) View() string {
	// The header
	s := state.question + "\n\n"

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
