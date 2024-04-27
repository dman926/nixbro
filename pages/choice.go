package pages

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Choice struct{}

func (c Choice) Update(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.page = WelcomePage
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor.Y > 0 {
				m.cursor.Up()
			}
		case "down", "j":
			if m.cursor.Y < len(m.choices)-1 {
				m.cursor.Down()
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor.Y]
			if ok {
				delete(m.selected, m.cursor.Y)
			} else {
				m.selected[m.cursor.Y] = struct{}{}
			}
		}
	}

	return m, nil
}

func (c Choice) View(m Model) string {
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor.Y == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress 1 to go to the other screen.\n"

	s += "\nPress q to quit.\n"

	return s
}