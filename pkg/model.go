package gomus

import (
	"errors"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faiface/beep"
)

var (
	termWidth  = 0
	termHeight = 0
)

type ModelConfig struct {
	MusicPath string
	GomusPath string
}

type Model struct {
	cursor           int
	currentlyPlaying int

	TrackPlayer
	TrackPlayerEffects
	trackIndex
	trackPlayerView
}

func NewModel(cfg ModelConfig) Model {
	if _, err := os.Stat(cfg.GomusPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(cfg.GomusPath, 0755)
		check(err)
	}

	return Model{
		cursor:           0,
		currentlyPlaying: 0,

		trackIndex:         NewDirTrackIndex(cfg),
		TrackPlayer:        NewTrackPlayer(),
		TrackPlayerEffects: NewTrackPlayerEffects(),
		trackPlayerView:    NewTrackPlayerView(),
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, tea.EnterAltScreen)
	cmds = append(cmds, m.trackPlayerView.Init())
	cmds = append(cmds, newLibraryUpdateCmd(m.trackIndex.tracks))
	return tea.Batch(cmds...)
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
			return m, tea.Quit
		case "enter":
			t := m.trackPlayerView.trackList.SelectedItem().(track)

			stream, format, err := t.GetStream()
			check(err)

			resampled := beep.Resample(4, format.SampleRate, m.TrackPlayer.SampleRate, stream)
			m.TrackPlayer.Play(resampled)

			cmds = append(cmds, newTrackChangeCmd(t))
		case " ":
			s := m.TrackPlayer.TogglePause()
			cmds = append(cmds, newTrackPauseCmd(s))
		case "-", "=":
			pe := &m.TrackPlayerEffects
			if msg.String() == "-" {
				pe.volume -= 0.1
				if pe.volume < minVolume {
					pe.volume = minVolume
				}
			} else {
				pe.volume += 0.1
				if pe.volume > maxVolume {
					pe.volume = maxVolume
				}
			}

			m.TrackPlayer.SetVolume(pe.volume)

			cmds = append(cmds, newTrackVolumeCmd(m.TrackPlayerEffects.volume))
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
