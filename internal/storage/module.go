package storage

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/auth"
	"github.com/smarthow/azure-for-dummies/internal/provider"
)

type FetchDoneMsg struct {
	Accounts []provider.StorageAccount
	Err      error
}

type Module struct {
	provider provider.StorageProvider
}

func New(ctx *auth.Context) *Module {
	return &Module{
		provider: newAzureProvider(ctx),
	}
}

func (m *Module) Name() string { return "Storage" }
func (m *Module) Icon() string { return "[ST]" }

func (m *Module) ListView() tea.Model {
	return newListView(m.provider)
}

func (m *Module) DetailView(id string) tea.Model {
	return newDetailView(id, m.provider)
}
