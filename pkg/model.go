package gomus

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	termWidth  = 0
	termHeight = 0
)

type ModelArgs struct {
	MusicPath string
}

type Model struct {
	cursor           int
	currentlyPlaying int

	trackPlayer
	trackIndex

	trackPlayerView trackPlayerView
}

func NewModel(args ModelArgs) Model {
	ti := NewDirTrackIndex(args.MusicPath)
	tpv := newTrackPlayerView(ti.tracks)

	return Model{
		cursor:           0,
		currentlyPlaying: 0,

		trackIndex:  ti,
		trackPlayer: trackPlayer{},

		trackPlayerView: tpv,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.trackPlayerView.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		termHeight = msg.Height
		termWidth = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.trackPlayer.close()
			return m, tea.Quit
		case "enter":
			t := m.trackPlayerView.trackList.SelectedItem().(track)
			cmds = append(cmds, newTrackChangeCmd(t))
			m.trackPlayer.play(t.getReader())
		}
	}
	var cmd tea.Cmd
	m.trackPlayerView, cmd = m.trackPlayerView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.trackPlayerView.View()
}
