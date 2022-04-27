package gomus

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type trackListDelegate struct{}

func (d trackListDelegate) Height() int                               { return 1 }
func (d trackListDelegate) Spacing() int                              { return 0 }
func (d trackListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d trackListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t := listItem.(track)
	f := t.fullName()

	if m.Index() == index {
		li := lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("[>] %s", f))
		fmt.Fprintf(w, li)
		return
	}

	fmt.Fprintf(w, fmt.Sprintf("    %s", f))
}

func newTrackListDelegate() trackListDelegate {
	return trackListDelegate{}
}
