package storage

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/auth"
	"github.com/smarthow/azure-for-dummies/internal/azutil"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/router"
	"github.com/smarthow/azure-for-dummies/internal/styles"
	"github.com/smarthow/azure-for-dummies/internal/tabs"
)

type detailFetchDoneMsg struct {
	account    provider.StorageAccount
	containers []provider.BlobContainer
	err        error
}

type detailView struct {
	id       string
	rg       string
	name     string
	provider *azureProvider
	auth     *auth.Context
	tabs     tabs.Model
	spinner  spinner.Model
	loading  bool
	width    int
	height   int
}

func newDetailView(id string, p *azureProvider, a *auth.Context) detailView {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = s.Style.Foreground(styles.Mauve)

	t := tabs.New(
		tabs.Tab{Title: "Overview", Content: newOverviewTab()},
		tabs.Tab{Title: "Containers", Content: newContainersTab()},
	)

	return detailView{
		id:       id,
		rg:       azutil.ExtractRG(id),
		name:     azutil.ExtractName(id),
		provider: p,
		auth:     a,
		tabs:     t,
		spinner:  s,
		loading:  true,
	}
}

func (v detailView) Init() tea.Cmd {
	return tea.Batch(v.spinner.Tick, v.fetchDetail())
}

func (v detailView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = msg.Width
		v.height = msg.Height
		contentMsg := tea.WindowSizeMsg{Width: msg.Width, Height: msg.Height - 1}
		var cmd tea.Cmd
		v.tabs, cmd = v.tabs.Update(contentMsg)
		return v, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			return v, func() tea.Msg { return router.PopMsg{} }
		}

	case detailFetchDoneMsg:
		v.loading = false
		if msg.err != nil {
			return v, nil
		}
		overview := newOverviewTab().SetAccount(msg.account)
		v.tabs = v.tabs.SetContent(0, overview)
		containers := newContainersTab().SetContainers(msg.containers)
		v.tabs = v.tabs.SetContent(1, containers)
		return v, nil
	}

	if v.loading {
		var cmd tea.Cmd
		v.spinner, cmd = v.spinner.Update(msg)
		return v, cmd
	}

	var cmd tea.Cmd
	v.tabs, cmd = v.tabs.Update(msg)
	return v, cmd
}

func (v detailView) View() string {
	if v.loading {
		return v.spinner.View() + " Loading " + v.name + "..."
	}
	return v.tabs.View()
}

func (v detailView) fetchDetail() tea.Cmd {
	rg, name := v.rg, v.name
	p := v.provider
	return func() tea.Msg {
		ctx := context.Background()
		account, err := p.GetStorageAccount(ctx, rg, name)
		if err != nil {
			return detailFetchDoneMsg{err: err}
		}
		containers, err := p.ListBlobContainers(ctx, rg, name)
		if err != nil {
			return detailFetchDoneMsg{account: account, err: err}
		}
		return detailFetchDoneMsg{account: account, containers: containers}
	}
}
