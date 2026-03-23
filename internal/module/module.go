package module

import tea "github.com/charmbracelet/bubbletea"

type Module interface {
	Name() string
	Icon() string
	ListView() tea.Model
	DetailView(id string) tea.Model
}
