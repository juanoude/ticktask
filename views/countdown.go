package views

import (
	"fmt"
	"os"
	"ticktask/player"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/common-nighthawk/go-figure"
	"github.com/ebitengine/oto/v3"
)

type UserStatus int

const (
	IdleStatus UserStatus = iota
	FocusStatus
)

type countdownState struct {
	totalTime       time.Duration
	restTime        time.Duration
	currentTime     time.Duration
	focusPlayer     *oto.Player
	focusController chan player.PlayerCommand
	restPlayer      *oto.Player
	restController  chan player.PlayerCommand
	status          UserStatus
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
		totalTime:       t,
		restTime:        0 * time.Second,
		currentTime:     current,
		focusPlayer:     player.InitFocusPlayer(),
		focusController: make(chan player.PlayerCommand),
		restPlayer:      player.InitRestPlayer(),
		restController:  make(chan player.PlayerCommand),
		status:          FocusStatus,
	}
}

type TickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (state countdownState) Init() tea.Cmd {
	go player.InitPlayerListener(state.focusPlayer, state.focusController)
	go player.InitPlayerListener(state.restPlayer, state.restController)
	state.focusController <- player.PlayCommand
	return tickEvery()
}

func (state countdownState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if state.currentTime >= state.totalTime {
		state.restController <- player.CloseCommand
		state.focusController <- player.CloseCommand
		return state, tea.Quit
	}

	switch msg := msg.(type) {
	case TickMsg:
		if state.status == FocusStatus {
			state.currentTime = state.currentTime + time.Second
		} else {
			state.restTime = state.restTime + time.Second
		}
		// Return your Every command again to loop.
		return state, tickEvery()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			state.restController <- player.CloseCommand
			state.focusController <- player.CloseCommand
			close(state.restController)
			close(state.focusController)
			return state, tea.Quit
		case " ":
			if state.status == FocusStatus {
				state.status = IdleStatus
				state.focusController <- player.PauseCommand
				state.restController <- player.PlayCommand
			} else {
				state.status = FocusStatus
				state.focusController <- player.PlayCommand
				state.restController <- player.PauseCommand
			}

			return state, nil
		}
	}

	return state, nil
}

func (state countdownState) View() string {
	focusFigure := figure.NewColorFigure(state.currentTime.String(), "", "green", true)
	restFigure := figure.NewColorFigure(state.restTime.String(), "", "red", true)
	totalView := focusFigure.ColorString() + "\n\n" + restFigure.ColorString()
	return totalView
}
