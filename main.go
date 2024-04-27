package main

import (
	"log"

	"github.com/dman926/nixbro/pages"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(pages.Start(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
