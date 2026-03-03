package storage

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/auth"
	"github.com/smarthow/azure-for-dummies/internal/provider"
)

type FetchDoneMsg struct {
	Accounts []provider.StorageAccount
	Err      error
}

type Module struct {
	auth     *auth.Context
	provider *azureProvider
}

func New(ctx *auth.Context) *Module {
	return &Module{
		auth:     ctx,
		provider: newAzureProvider(ctx),
	}
}

func (m *Module) Name() string { return "Storage" }
func (m *Module) Icon() string { return "[ST]" }

func (m *Module) Fetch(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		accounts, err := m.provider.ListStorageAccounts(ctx, m.auth.SubscriptionID)
		return FetchDoneMsg{Accounts: accounts, Err: err}
	}
}

func (m *Module) ListView() tea.Model {
	return newListView(m.provider, m.auth)
}

func (m *Module) DetailView(id string) tea.Model {
	return newDetailView(id, m.provider, m.auth)
}
