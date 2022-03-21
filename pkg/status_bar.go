package gomus

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)
)

type statusBar struct {
	isPaused     bool
	currentTrack track
}

func (s statusBar) Init() tea.Cmd {
	return nil
}

func (s statusBar) View() string {
	w := lipgloss.Width

	var statusKey string
	if s.isPaused {
		statusKey = statusStyle.Render("PAUSED")
	} else {
		statusKey = statusStyle.Render("NOW PLAYING")
	}

	statusVal := statusText.Copy().
		Width(termWidth - w(statusKey)).
		Render(s.currentTrack.fullName())

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
	)

	return statusBarStyle.Width(termWidth).Render(bar)
}

func (s statusBar) Update(msg tea.Msg) (statusBar, tea.Cmd) {
	switch msg := msg.(type) {
	case trackChangeMsg:
		s.currentTrack = msg.nextTrack
	case trackPauseMsg:
		s.isPaused = msg.isPaused
	}

	return s, nil
}
