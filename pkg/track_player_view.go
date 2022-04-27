package gomus

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type trackPlayerView struct {
	trackList list.Model
	statusBar
}

func NewTrackPlayerView() trackPlayerView {
	l := list.New([]list.Item{}, newTrackListDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	s := statusBar{currentVolume: startVolume}
	return trackPlayerView{statusBar: s, trackList: l}
}

func (v trackPlayerView) Init() tea.Cmd {
	return nil
}

func (v trackPlayerView) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, v.trackList.View(), v.statusBar.View())
}

func (v trackPlayerView) Update(msg tea.Msg) (trackPlayerView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.trackList.SetHeight(msg.Height - 2)
		v.trackList.SetWidth(msg.Width)
	case libraryUpdateMsg:
		c := MapList(msg.tracks, func(t track) list.Item {
			return t
		})
		v.trackList.SetItems(c)
	}

	var cmd tea.Cmd
	v.trackList, cmd = v.trackList.Update(msg)
	cmds = append(cmds, cmd)

	v.statusBar, cmd = v.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
}
