package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smarthow/azure-for-dummies/internal/auth"
	"github.com/smarthow/azure-for-dummies/internal/containerapps"
	"github.com/smarthow/azure-for-dummies/internal/home"
	"github.com/smarthow/azure-for-dummies/internal/loading"
	"github.com/smarthow/azure-for-dummies/internal/module"
	"github.com/smarthow/azure-for-dummies/internal/postgres"
	"github.com/smarthow/azure-for-dummies/internal/router"
	"github.com/smarthow/azure-for-dummies/internal/storage"
	"github.com/smarthow/azure-for-dummies/internal/statusbar"
	"github.com/smarthow/azure-for-dummies/internal/styles"
	"github.com/smarthow/azure-for-dummies/internal/subscriptionpicker"
)

type Model struct {
	router    router.Model
	statusbar statusbar.Model
	registry  *module.Registry
	width     int
	height    int
}

func New() Model {
	return Model{
		router:   router.New(loading.New()),
		statusbar: statusbar.New(),
		registry: module.NewRegistry(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.router.Init(), auth.ListSubscriptions())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusbar = m.statusbar.WithWidth(msg.Width)

	case auth.SubscriptionsMsg:
		if len(msg.Subscriptions) == 1 {
			return m, auth.ResolveWithSubscription(msg.Subscriptions[0].ID)
		}
		var cmd tea.Cmd
		m.router, cmd = m.router.ReplaceRoot(subscriptionpicker.New(msg.Subscriptions))
		return m, cmd

	case auth.AuthReadyMsg:
		ctx := msg.Ctx
		m.registry.Register(containerapps.New(ctx))
		m.registry.Register(postgres.New(ctx))
		m.registry.Register(storage.New(ctx))
		m.statusbar = m.statusbar.
			WithSubscription(ctx.SubscriptionName).
			WithUser(ctx.UserPrincipal)
		var cmd tea.Cmd
		m.router, cmd = m.router.ReplaceRoot(home.New(m.registry.All()))
		return m, cmd

	case auth.AuthErrMsg:
		m.statusbar = m.statusbar.WithError(msg.Err)
		return m, nil
	}

	var cmd tea.Cmd
	m.router, cmd = m.router.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	statusHeight := 1
	contentHeight := m.height - statusHeight
	if contentHeight < 0 {
		contentHeight = 0
	}
	content := styles.ContentArea.
		Width(m.width).
		Height(contentHeight).
		Render(m.router.View())
	return lipgloss.JoinVertical(lipgloss.Left, content, m.statusbar.View())
}
