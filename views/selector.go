// Package views provides terminal UI components using the Bubble Tea framework.
// It includes interactive selectors, text inputs, and the focus timer countdown display.
package views

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// selectorState holds the state for the interactive list selector.
// Implements the tea.Model interface for Bubble Tea.
type selectorState struct {
	question string   // The prompt displayed above the list
	tasks    []string // The selectable options
	cursor   int      // Current cursor position (0-indexed)
	Selected int      // The index of the selected item (-1 if cancelled)
}

// RunSelector displays an interactive list selector and returns the selected index.
// Returns -1 if the user cancels (presses 'q' or Ctrl+C).
// Navigation: up/k = move up, down/j = move down, enter/space = select, q = cancel.
func RunSelector(options []string, question string) int {
	p := tea.NewProgram(initSelector(options, question))
	result, err := p.Run()
	if err != nil {
		fmt.Printf("Oops, there's been an error: %v", err)
		os.Exit(1)
	}

	finalResult := result.(selectorState)
	return finalResult.Selected
}

// initSelector creates the initial state for a selector with the given options and prompt.
func initSelector(options []string, questionString string) selectorState {
	return selectorState{
		question: questionString,
		tasks:    options,
		cursor:   0,
		Selected: 0,
	}
}

// Init is called when the program starts. Returns nil as no initial I/O is needed.
func (state selectorState) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and updates the selector state.
// Supported keys:
//   - up/k: Move cursor up
//   - down/j: Move cursor down
//   - enter/space: Select current item and exit
//   - q/ctrl+c: Cancel and exit (sets Selected to -1)
func (state selectorState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			state.Selected = -1
			return state, tea.Quit

		case "up", "k":
			if state.cursor > 0 {
				state.cursor--
			}

		case "down", "j":
			if state.cursor < len(state.tasks)-1 {
				state.cursor++
			}

		case "enter", " ":
			state.Selected = state.cursor
			return state, tea.Quit
		}
	}

	return state, nil
}

// View renders the selector UI as a string.
// Shows the question prompt, list of options with cursor indicator (->),
// and footer instructions.
func (state selectorState) View() string {
	s := state.question + "\n\n"

	for i, choice := range state.tasks {
		cursor := " "
		if state.cursor == i {
			cursor = "->"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n"
	return s
}
