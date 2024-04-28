package pages

import (
	"fmt"
	"io"

	"github.com/dman926/nixbro/cursor"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const APP_NAME string = "Nix Bro"

type PageModel interface {
	Page() Page
	PageTitle() pageTitle

	View(m Model) string
	Update(m *Model, msg tea.Msg) (tea.Model, tea.Cmd)
}

type Page int

type Model struct {
	cursor cursor.Cursor
	page   Page
	pages  []PageModel

	navigatorOpen bool
	navigator     list.Model
}

const (
	WelcomePage Page = iota
	InstallPage
)

var pages = []PageModel{WelcomeInit(), InstallInit()}

/* TODO Move somewhere. Probably generalize. pageTitle becomes Page. */

func underline(s string) string {
	return "\033[4m" + s + "\033[0m"
}

type pageTitle string

func (i pageTitle) FilterValue() string { return "" }

type pageTitleDelegate struct{}

func (pageTitleDelegate) Height() int                               { return 1 }
func (pageTitleDelegate) Spacing() int                              { return 0 }
func (pageTitleDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (pageTitleDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(pageTitle)
	if !ok {
		return
	}

	cursor := "  "
	if m.Index() == index {
		cursor = "> "
	}

	str := fmt.Sprintf("%s %s", cursor, i)
	if m.Index() == index {
		str = underline(str)
	}

	fmt.Fprintln(w, str)
}

/* End TODO notice */

func Start() Model {
	navigatorItems := make([]list.Item, len(pages))
	for i, page := range pages {
		navigatorItems[i] = page.PageTitle()
	}

	navigator := list.New(navigatorItems, pageTitleDelegate{}, 20, 14)

	return Model{
		cursor:    cursor.New(),
		page:      WelcomePage,
		pages:     pages,
		navigator: navigator,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle(APP_NAME)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+n":
			m.navigatorOpen = !m.navigatorOpen
		case "enter":
			if m.navigatorOpen {
				i, ok := m.navigator.SelectedItem().(pageTitle)
				if ok {
					for _, page := range m.pages {
						if page.PageTitle() == i {
							m.page = page.Page()
							m.navigatorOpen = false
							return m, nil
						}
					}
				}
			}

			return m._update(msg)
		default:
			return m._update(msg)
		}
	}

	return m._update(msg)
}

func (m Model) _update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.navigatorOpen {
		var cmd tea.Cmd
		m.navigator, cmd = m.navigator.Update(msg)
		return m, cmd
	}

	return m.currentPage().Update(&m, msg)

}

func (m Model) View() string {
	if m.navigatorOpen {
		return "\n" + m.navigator.View()
	}

	s := m.currentPage().View(m)

	s += "\nPress CTRL+N to toggle the navigator.\n"

	s += "Press q to quit.\n"

	return s
}

func (m Model) currentPage() PageModel {
	return m.pages[m.page]
}
