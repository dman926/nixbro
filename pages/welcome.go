package pages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Welcome struct{}

func (w Welcome) Update(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.page = ChoicePage
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (w Welcome) View(m Model) string {
	s := "This is the other screen\n"

	s += "\nPress 1 to go to the other screen.\n"

	s += "\nPress q to quit.\n"

	return s
}
