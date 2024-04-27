package pages

import (
	"github.com/dman926/nixbro/cursor"

	tea "github.com/charmbracelet/bubbletea"
)

type PageModel interface {
	View(m Model) string
	Update(m Model, msg tea.Msg) (tea.Model, tea.Cmd)
}

type Page int

type Model struct {
	cursor cursor.Cursor
	page   Page
	pages  []PageModel
}

const (
	WelcomePage Page = iota
	ChoicePage
)

func Start() Model {
	return Model{
		cursor: cursor.New(),
		page:   WelcomePage,
		pages:  []PageModel{WelcomeInit(), ChoiceInit()},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("Nix Bro")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.pages[m.page].Update(m, msg)
}

func (m Model) View() string {
	return m.pages[m.page].View(m)
}
