package containerapps

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/auth"
	"github.com/smarthow/azure-for-dummies/internal/provider"
)

type FetchDoneMsg struct {
	Apps []provider.ContainerApp
	Err  error
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

func (m *Module) Name() string { return "Container Apps" }
func (m *Module) Icon() string { return "[CA]" }

func (m *Module) Fetch(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		apps, err := m.provider.ListContainerApps(ctx, m.auth.SubscriptionID)
		return FetchDoneMsg{Apps: apps, Err: err}
	}
}

func (m *Module) ListView() tea.Model {
	return newListView(m.provider, m.auth)
}

func (m *Module) DetailView(id string) tea.Model {
	return newDetailView(id, m.provider, m.auth)
}
