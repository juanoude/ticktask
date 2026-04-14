package views

import (
	"fmt"
	"os"
	"ticktask/player"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/common-nighthawk/go-figure"
)

// UserStatus represents the current timer mode.
type UserStatus int

const (
	// IdleStatus indicates rest/break mode (5-minute default).
	IdleStatus UserStatus = iota
	// FocusStatus indicates focused work mode (25-minute default Pomodoro).
	FocusStatus
	// GenericStatus indicates general/chore mode (2-hour limit).
	GenericStatus
)

// countdownState holds the state for the Pomodoro-style focus timer.
// Manages three separate timers (focus, rest, generic) with corresponding audio players.
// Implements the tea.Model interface for Bubble Tea.
type countdownState struct {
	rest    *countdown // Rest/break timer (5 min default)
	focus   *countdown // Focus work timer (25 min default)
	generic *countdown // Generic/chore timer (2 hour limit)
	status  UserStatus // Current active timer mode
	isOpen  bool       // If true, timers don't auto-rotate (open-ended mode)
}

// countdown represents a single timer with associated audio player.
type countdown struct {
	currentTime time.Duration    // Time elapsed in current session
	player      *player.TTPlayer // Audio player for this timer mode
	totalTime   time.Duration    // Total time accumulated across sessions
	limit       time.Duration    // Duration before auto-rotating to next mode
}

// RunCountdown starts the focus timer interface.
// If isOpenFlag is true, the timer runs indefinitely without auto-rotating.
// Otherwise, follows Pomodoro pattern: 25min focus → 5min rest → repeat.
func RunCountdown(isOpenFlag bool) {
	p := tea.NewProgram(initCountdown(isOpenFlag))
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Oops, there's been an error: %v", err)
		os.Exit(1)
	}
}

// initCountdown creates the initial timer state with three players.
// Loads music for each mode (focus, rest, generic) from configured backend.
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

// TickMsg is sent every second to update the timer display.
type TickMsg time.Time

// tickEvery returns a command that sends TickMsg every second.
func tickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Init initializes all three audio players and starts focus mode.
// Called when the Bubble Tea program starts.
func (state countdownState) Init() tea.Cmd {
	state.focus.player.InitPlayer()
	state.rest.player.InitPlayer()
	state.generic.player.InitPlayer()
	state.focus.player.Play()
	return tickEvery()
}

// Focus switches to focus mode: plays focus music, pauses others.
func (state *countdownState) Focus() {
	state.status = FocusStatus
	state.focus.player.Play()
	state.rest.player.Pause()
	state.generic.player.Pause()
}

// Rest switches to rest mode: plays rest music, pauses others.
func (state *countdownState) Rest() {
	state.status = IdleStatus
	state.focus.player.Pause()
	state.rest.player.Play()
	state.generic.player.Pause()
}

// Chore switches to generic/chore mode: plays generic music, pauses others.
func (state *countdownState) Chore() {
	state.status = GenericStatus
	state.focus.player.Pause()
	state.rest.player.Pause()
	state.generic.player.Play()
}

// Count increments the current session time by one second.
func (cd *countdown) Count() {
	cd.currentTime += time.Second
}

// Rotate resets current session time and adds the limit to total time.
// Called when auto-rotating between focus and rest modes.
func (cd *countdown) Rotate() {
	cd.currentTime = 0 * time.Second
	cd.totalTime += cd.limit
}

// TotalTime returns the total accumulated time including current session.
func (cd *countdown) TotalTime() time.Duration {
	return cd.totalTime + cd.currentTime
}

// Update handles timer ticks and keyboard input.
// In standard mode (!isOpen), auto-rotates between focus and rest when limits are reached.
// Keyboard controls:
//   - Space: Toggle between focus and rest modes
//   - Backspace: Toggle between current mode and chore mode
//   - q/Ctrl+C: Quit and close all players
func (state countdownState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Auto-rotate between focus and rest in standard Pomodoro mode
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
		// Increment the appropriate timer based on current mode
		switch state.status {
		case FocusStatus:
			state.focus.Count()
		case IdleStatus:
			state.rest.Count()
		case GenericStatus:
			state.generic.Count()
		}
		return state, tickEvery()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Clean up audio players before exiting
			state.rest.player.Close()
			state.focus.player.Close()
			state.generic.player.Close()
			return state, tea.Quit

		case " ":
			// Space toggles between focus and rest
			if state.status == FocusStatus {
				state.Rest()
			} else {
				state.Focus()
			}
			return state, nil

		case "backspace":
			// Backspace toggles chore mode
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

// View renders the timer display using ASCII art figures.
// Shows three timers in different colors:
//   - Green: Focus time
//   - Red: Rest time
//   - Blue: Chore/generic time
func (state countdownState) View() string {
	focusFigure := figure.NewColorFigure(state.focus.TotalTime().String(), "", "green", true)
	restFigure := figure.NewColorFigure(state.rest.TotalTime().String(), "", "red", true)
	genericFigure := figure.NewColorFigure(state.generic.TotalTime().String(), "", "blue", true)
	totalView := focusFigure.ColorString() + "\n\n" + restFigure.ColorString() + "\n\n" + genericFigure.ColorString()
	return totalView
}
