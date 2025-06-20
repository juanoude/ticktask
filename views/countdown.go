package views

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/common-nighthawk/go-figure"
)

type countdownState struct {
	totalTime   time.Duration
	currentTime time.Duration
}

func RunCountdown(time time.Duration) {
	p := tea.NewProgram(initCountdown(time))
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Oops, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initCountdown(t time.Duration) countdownState {
	current := 0 * time.Second
	return countdownState{
		totalTime:   t,
		currentTime: current,
	}
}

type TickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (state countdownState) Init() tea.Cmd {
	return tickEvery()
}

func (state countdownState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if state.currentTime >= state.totalTime {
		return state, tea.Quit
	}

	switch msg := msg.(type) {
	case TickMsg:
		state.currentTime = state.currentTime + time.Second
		// Return your Every command again to loop.
		return state, tickEvery()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return state, tea.Quit
		}
	}

	return state, nil
}

func (state countdownState) View() string {
	myFigure := figure.NewColorFigure(state.currentTime.String(), "", "red", true)
	return myFigure.ColorString()
}
