package tabs

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type TabChangedMsg struct{ Index int }

type Tab struct {
	Title   string
	Content tea.Model
}

type Model struct {
	tabs   []Tab
	active int
	width  int
	height int
}

func New(tabs ...Tab) Model {
	return Model{tabs: tabs}
}

func (m Model) Active() int      { return m.active }
func (m Model) ActiveTab() Tab   { return m.tabs[m.active] }
func (m Model) Tabs() []Tab      { return m.tabs }
func (m Model) Width() int       { return m.width }
func (m Model) ContentHeight() int { return m.height - 1 }

func (m Model) SetContent(idx int, content tea.Model) Model {
	m.tabs[idx].Content = content
	return m
}

func (m Model) Init() tea.Cmd {
	if len(m.tabs) == 0 {
		return nil
	}
	return m.tabs[m.active].Content.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		contentMsg := tea.WindowSizeMsg{Width: msg.Width, Height: m.ContentHeight()}
		var cmds []tea.Cmd
		for i, t := range m.tabs {
			updated, cmd := t.Content.Update(contentMsg)
			m.tabs[i].Content = updated
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			prev := m.active
			m.active = (m.active + 1) % len(m.tabs)
			if m.active != prev {
				return m, func() tea.Msg { return TabChangedMsg{Index: m.active} }
			}
			return m, nil
		case "shift+tab":
			prev := m.active
			m.active = (m.active - 1 + len(m.tabs)) % len(m.tabs)
			if m.active != prev {
				return m, func() tea.Msg { return TabChangedMsg{Index: m.active} }
			}
			return m, nil
		default:
			if n, err := strconv.Atoi(msg.String()); err == nil && n >= 1 && n <= len(m.tabs) {
				prev := m.active
				m.active = n - 1
				if m.active != prev {
					return m, func() tea.Msg { return TabChangedMsg{Index: m.active} }
				}
				return m, nil
			}
		}
	}

	if len(m.tabs) == 0 {
		return m, nil
	}
	updated, cmd := m.tabs[m.active].Content.Update(msg)
	m.tabs[m.active].Content = updated
	return m, cmd
}

func (m Model) View() string {
	if len(m.tabs) == 0 {
		return ""
	}

	var tabHeaders []string
	for i, t := range m.tabs {
		if i == m.active {
			tabHeaders = append(tabHeaders, styles.TabActive.Render(t.Title))
		} else {
			tabHeaders = append(tabHeaders, styles.TabInactive.Render(t.Title))
		}
	}

	bar := styles.TabBar.Width(m.width).Render(strings.Join(tabHeaders, ""))
	content := m.tabs[m.active].Content.View()
	return bar + "\n" + content
}
