package gomus

import tea "github.com/charmbracelet/bubbletea"

type trackChangeMsg struct {
	nextTrack track
}

func newTrackChangeCmd(nextTrack track) tea.Cmd {
	return func() tea.Msg {
		return trackChangeMsg{nextTrack}
	}
}
