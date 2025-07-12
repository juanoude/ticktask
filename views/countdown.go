package views

import (
	"fmt"
	"os"
	"ticktask/player"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/common-nighthawk/go-figure"
)

type UserStatus int

const (
	IdleStatus UserStatus = iota
	FocusStatus
	GenericStatus
)

type countdownState struct {
	rest    *countdown
	focus   *countdown
	generic *countdown
	status  UserStatus
	isOpen  bool
}

type countdown struct {
	currentTime time.Duration
	player      *player.TTPlayer
	totalTime   time.Duration
	limit       time.Duration
}

func RunCountdown(isOpenFlag bool) {
	p := tea.NewProgram(initCountdown(isOpenFlag))
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Oops, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initCountdown(isOpenFlag bool) countdownState {
	return countdownState{
		status: FocusStatus,
		isOpen: isOpenFlag,
		focus: &countdown{
			player:      player.GetFocusPlayer(),
			totalTime:   0 * time.Second,
			currentTime: 0 * time.Second,
			limit:       25 * time.Minute,
		},
		rest: &countdown{
			player:      player.GetRestPlayer(),
			totalTime:   0 * time.Second,
			currentTime: 0 * time.Second,
			limit:       5 * time.Minute,
		},
		generic: &countdown{
			player:      player.GetGenericPlayer(),
			totalTime:   0 * time.Second,
			currentTime: 0 * time.Second,
			limit:       2 * time.Hour,
		},
	}
}

type TickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (state countdownState) Init() tea.Cmd {
	state.focus.player.InitPlayer()
	state.rest.player.InitPlayer()
	state.generic.player.InitPlayer()
	state.focus.player.Play()
	return tickEvery()
}

func (state *countdownState) Focus() {
	state.status = FocusStatus
	state.focus.player.Play()
	state.rest.player.Pause()
	state.generic.player.Pause()
}

func (state *countdownState) Rest() {
	state.status = IdleStatus
	state.focus.player.Pause()
	state.rest.player.Play()
	state.generic.player.Pause()
}

func (state *countdownState) Chore() {
	state.status = GenericStatus
	state.focus.player.Pause()
	state.rest.player.Pause()
	state.generic.player.Play()
}

func (cd *countdown) Count() {
	cd.currentTime += time.Second
}

func (cd *countdown) Rotate() {
	cd.currentTime = 0 * time.Second
	cd.totalTime += cd.limit
}

func (cd *countdown) TotalTime() time.Duration {
	return cd.totalTime + cd.currentTime
}

func (state countdownState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !state.isOpen {
		if state.focus.currentTime >= state.focus.limit {
			state.focus.Rotate()
			state.Rest()
			return state, tickEvery()
		}

		if state.rest.currentTime >= state.rest.limit {
			state.rest.Rotate()
			state.Focus()
			return state, tickEvery()
		}
	}

	switch msg := msg.(type) {
	case TickMsg:
		switch state.status {
		case FocusStatus:
			state.focus.Count()
		case IdleStatus:
			state.rest.Count()
		case GenericStatus:
			state.generic.Count()
		}
		// Return your Every command again to loop.
		return state, tickEvery()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			state.rest.player.Close()
			state.focus.player.Close()
			state.generic.player.Close()
			return state, tea.Quit
		case " ":
			if state.status == FocusStatus {
				state.Rest()
			} else {
				state.Focus()
			}

			return state, nil
		case "backspace":
			if state.status == GenericStatus {
				state.Rest()
			} else {
				state.Chore()
			}

			return state, nil
		}
	}

	return state, nil
}

func (state countdownState) View() string {
	focusFigure := figure.NewColorFigure(state.focus.TotalTime().String(), "", "green", true)
	restFigure := figure.NewColorFigure(state.rest.TotalTime().String(), "", "red", true)
	genericFigure := figure.NewColorFigure(state.generic.TotalTime().String(), "", "blue", true)
	totalView := focusFigure.ColorString() + "\n\n" + restFigure.ColorString() + "\n\n" + genericFigure.ColorString()
	return totalView
}
