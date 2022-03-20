package main

import (
	"log"

	gomus "git.cesium.pw/niku/gomus/pkg"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(gomus.NewModel())
	if err := p.Start(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
