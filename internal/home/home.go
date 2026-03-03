package home

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smarthow/azure-for-dummies/internal/auth"
	"github.com/smarthow/azure-for-dummies/internal/module"
	"github.com/smarthow/azure-for-dummies/internal/router"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type Model struct {
	modules []module.Module
	cursor  int
	width   int
	height  int
	auth    *auth.Context
}

func New(modules []module.Module, ctx *auth.Context) Model {
	return Model{modules: modules, auth: ctx}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}
		case "right", "l":
			if m.cursor < len(m.modules)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.modules) == 0 {
				break
			}
			mod := m.modules[m.cursor]
			return m, tea.Batch(
				func() tea.Msg { return router.PushMsg{Screen: mod.ListView()} },
				mod.Fetch(context.Background()),
			)
		}
	}
	return m, nil
}

func (m Model) View() string {
	if len(m.modules) == 0 {
		return styles.Placeholder.Render("No modules registered.")
	}
	cards := make([]string, len(m.modules))
	for i, mod := range m.modules {
		label := mod.Icon() + " " + mod.Name()
		if i == m.cursor {
			cards[i] = styles.ModuleCardSelected.Render(label)
		} else {
			cards[i] = styles.ModuleCard.Render(label)
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, cards...)
}
