package subscriptionpicker

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/auth"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type subItem struct {
	sub auth.Subscription
}

func (i subItem) Title() string       { return i.sub.Name }
func (i subItem) Description() string { return i.sub.ID }
func (i subItem) FilterValue() string { return i.sub.Name }

type Model struct {
	list list.Model
}

func New(subs []auth.Subscription) Model {
	items := make([]list.Item, len(subs))
	for i, s := range subs {
		items[i] = subItem{sub: s}
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Subscription"
	l.Styles.Title = styles.ListTitle
	return Model{list: l}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if sel, ok := m.list.SelectedItem().(subItem); ok {
				return m, auth.ResolveWithSubscription(sel.sub.ID)
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
