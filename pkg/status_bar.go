package gomus

import (
	"fmt"

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
			Padding(0, 1)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle).Padding(0, 1)
)

type statusBar struct {
	isPaused      bool
	currentTrack  track
	currentVolume float64
}

func (s *statusBar) Init() tea.Cmd {
	return nil
}

func (s statusBar) View() string {
	w := lipgloss.Width

	var statusKey string
	if s.isPaused {
		statusKey = statusStyle.Render("")
	} else {
		statusKey = statusStyle.Render("")
	}

	v := MapFloatBetween(s.currentVolume, minVolume, maxVolume, 0, 100)
	vs := fmt.Sprintf("vol %d", int(v))

	statusVal := statusText.Copy().
		Width(termWidth - w(statusKey) - w(vs) - 2).
		Render(s.currentTrack.fullName())

	statusVol := statusStyle.Copy().Align(lipgloss.Right).Width(w(vs) + 2).Render(vs)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		statusVol,
	)

	return statusBarStyle.Render(bar)
}

func (s statusBar) Update(msg tea.Msg) (statusBar, tea.Cmd) {
	switch msg := msg.(type) {
	case trackChangeMsg:
		s.currentTrack = msg.nextTrack
	case trackPauseMsg:
		s.isPaused = msg.isPaused
	case trackVolumeMsg:
		s.currentVolume = msg.volume
	}

	return s, nil
}
