package views

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errorMsg error
)

type inputState struct {
	header       string
	textInput    textinput.Model
	wasCancelled bool
	isSecret     bool
	err          error
}

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

func (state inputState) Init() tea.Cmd {
	return textinput.Blink
}

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

	// We handle errors just like any other message
	case errorMsg:
		state.err = msg
		return state, nil
	}

	state.textInput, cmd = state.textInput.Update(msg)
	return state, cmd

}

func (state inputState) View() string {
	return fmt.Sprintf(
		state.header+"\n\n%s\n\n%s",
		state.textInput.View(),
		"(esc/ctrl+c to quit)\n",
		"(enter to confirm)",
	) + "\n"
}
