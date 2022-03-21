package main

import (
	"log"
	"os"

	gomus "git.cesium.pw/niku/gomus/pkg"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected a path to some music")
	}

	p := tea.NewProgram(gomus.NewModel(gomus.ModelArgs{
		MusicPath: os.Args[1],
	}))

	if err := p.Start(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
