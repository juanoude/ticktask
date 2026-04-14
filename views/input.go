package views

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// errorMsg wraps errors for the Bubble Tea message system.
type (
	errorMsg error
)

// inputState holds the state for the text input component.
// Implements the tea.Model interface for Bubble Tea.
type inputState struct {
	header       string          // Prompt displayed above the input field
	textInput    textinput.Model // The Bubble Tea text input component
	wasCancelled bool            // True if user pressed Esc or Ctrl+C
	isSecret     bool            // If true, input should be masked (not currently implemented)
	err          error           // Any error that occurred during input
}

// RunInput displays a text input prompt and returns the entered value.
// Returns the input text and a boolean indicating if the user cancelled.
// The isSecret parameter is intended for password masking but not currently implemented.
func RunInput(header string, isSecret bool) (string, bool) {
	p := tea.NewProgram(initInput(header, isSecret))
	result, err := p.Run()
	if err != nil {
		log.Printf("Oops, error happened during input: %v", err)
		os.Exit(1)
	}

	finalResult := result.(inputState)
	return finalResult.textInput.Value(), finalResult.wasCancelled
}

// initInput creates the initial state for a text input with the given header.
// Sets up a text input with 156 character limit and 20 character display width.
func initInput(header string, isSecret bool) inputState {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return inputState{
		header:       header,
		textInput:    ti,
		wasCancelled: false,
		isSecret:     false,
		err:          nil,
	}
}

// Init starts the cursor blinking animation.
func (state inputState) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles keyboard input for the text field.
// Enter confirms the input, Esc/Ctrl+C cancels.
func (state inputState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return state, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			state.wasCancelled = true
			return state, tea.Quit
		}

	case errorMsg:
		state.err = msg
		return state, nil
	}

	state.textInput, cmd = state.textInput.Update(msg)
	return state, cmd

}

// View renders the input prompt with the header, text field, and instructions.
func (state inputState) View() string {
	return fmt.Sprintf(
		state.header+"\n\n%s\n\n%s",
		state.textInput.View(),
		"(esc/ctrl+c to quit)\n",
		"(enter to confirm)",
	) + "\n"
}
