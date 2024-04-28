package pages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Install struct{}

func InstallInit() Install {
	return Install{}
}

func (i Install) Page() Page {
	return InstallPage
}

func (i Install) PageTitle() pageTitle {
	return pageTitle("Install")
}

func (i Install) Update(m *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (i Install) View(m Model) string {
	s := "This is the install page\n"

	return s
}
