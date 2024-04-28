package pages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Welcome struct{}

func WelcomeInit() Welcome {
	return Welcome{}
}

func (w Welcome) Page() Page {
	return WelcomePage
}

func (w Welcome) PageTitle() pageTitle {
	return pageTitle("Welcome")
}

func (w Welcome) Update(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (w Welcome) View(m Model) string {
	s := "This is the welcome page\n"

	return s
}
