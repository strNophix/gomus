package gomus

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faiface/beep/speaker"
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

	TrackPlayer
	trackIndex

	trackPlayerView
}

func NewModel(args ModelArgs) Model {
	ti := NewDirTrackIndex(args.MusicPath)
	tpv := newTrackPlayerView(ti.tracks)

	return Model{
		cursor:           0,
		currentlyPlaying: 0,

		trackIndex:  ti,
		TrackPlayer: TrackPlayer{},

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
			m.TrackPlayer.Close()
			return m, tea.Quit
		case "enter":
			t := m.trackPlayerView.trackList.SelectedItem().(track)
			cmds = append(cmds, newTrackChangeCmd(t))
			stream, format, err := t.GetStream()
			check(err)
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			m.TrackPlayer.Play(&stream)
		case " ":
			pauseState := m.TrackPlayer.TogglePause()
			cmds = append(cmds, newTrackPauseCmd(pauseState))
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
