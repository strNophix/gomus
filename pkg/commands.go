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

type trackPauseMsg struct {
	isPaused bool
}

func newTrackPauseCmd(isPaused bool) tea.Cmd {
	return func() tea.Msg {
		return trackPauseMsg{isPaused}
	}
}

type trackVolumeMsg struct {
	volume float64
}

func newTrackVolumeCmd(volume float64) tea.Cmd {
	return func() tea.Msg {
		return trackVolumeMsg{volume}
	}
}
