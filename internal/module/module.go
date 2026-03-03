package module

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type Module interface {
	Name() string
	Icon() string
	Fetch(ctx context.Context) tea.Cmd
	ListView() tea.Model
	DetailView(id string) tea.Model
}
