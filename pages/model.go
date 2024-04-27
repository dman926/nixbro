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

	choices  []string
	selected map[int]struct{}
}

const (
	ChoicePage Page = iota
	WelcomePage
)

func Start() Model {
	return Model{
		page:  WelcomePage,
		pages: []PageModel{Choice{}, Welcome{}},

		cursor: cursor.New(),

		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("Grocery List")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.pages[m.page].Update(m, msg)
}

func (m Model) View() string {
	return m.pages[m.page].View(m)
}
