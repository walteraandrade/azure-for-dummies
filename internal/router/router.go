package router

import (
	tea "github.com/charmbracelet/bubbletea"
)

type PushMsg struct{ Screen tea.Model }
type PopMsg struct{}

type Model struct {
	stack  []tea.Model
	width  int
	height int
}

func New(initial tea.Model) Model {
	return Model{stack: []tea.Model{initial}}
}

func (m Model) Init() tea.Cmd {
	if len(m.stack) == 0 {
		return nil
	}
	return m.stack[len(m.stack)-1].Init()
}

func (m Model) sizeCmd() tea.Cmd {
	w, h := m.width, m.height
	if w == 0 && h == 0 {
		return nil
	}
	return func() tea.Msg { return tea.WindowSizeMsg{Width: w, Height: h} }
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case PushMsg:
		m.stack = push(m.stack, msg.Screen)
		return m, tea.Batch(msg.Screen.Init(), m.sizeCmd())
	case PopMsg:
		if len(m.stack) > 1 {
			m.stack = pop(m.stack)
			return m, m.sizeCmd()
		}
		return m, nil
	}

	if len(m.stack) == 0 {
		return m, nil
	}

	top := m.stack[len(m.stack)-1]
	updated, cmd := top.Update(msg)
	m.stack = replace(m.stack, updated)
	return m, cmd
}

func (m Model) View() string {
	if len(m.stack) == 0 {
		return ""
	}
	return m.stack[len(m.stack)-1].View()
}

func push(stack []tea.Model, screen tea.Model) []tea.Model {
	result := make([]tea.Model, len(stack)+1)
	copy(result, stack)
	result[len(stack)] = screen
	return result
}

func pop(stack []tea.Model) []tea.Model {
	result := make([]tea.Model, len(stack)-1)
	copy(result, stack[:len(stack)-1])
	return result
}

func replace(stack []tea.Model, top tea.Model) []tea.Model {
	result := make([]tea.Model, len(stack))
	copy(result, stack)
	result[len(stack)-1] = top
	return result
}

func (m Model) ReplaceRoot(screen tea.Model) (Model, tea.Cmd) {
	m.stack = []tea.Model{screen}
	return m, tea.Batch(screen.Init(), m.sizeCmd())
}
