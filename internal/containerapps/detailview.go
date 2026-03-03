package containerapps

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
	app  provider.ContainerApp
	revs []provider.RevisionInfo
	err  error
}

type logStreamReadyMsg struct {
	ch <-chan provider.LogEntry
}

type detailView struct {
	id         string
	rg         string
	name       string
	provider   *azureProvider
	auth       *auth.Context
	tabs       tabs.Model
	spinner    spinner.Model
	loading    bool
	logStarted bool
	logCancel  context.CancelFunc
	width      int
	height     int
}

func newDetailView(id string, p *azureProvider, a *auth.Context) detailView {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = s.Style.Foreground(styles.Mauve)

	t := tabs.New(
		tabs.Tab{Title: "Overview", Content: newOverviewTab()},
		tabs.Tab{Title: "Revisions", Content: newRevisionsTab()},
		tabs.Tab{Title: "Logs", Content: newLogsTab()},
		tabs.Tab{Title: "Settings", Content: newSettingsTab()},
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
			if v.logCancel != nil {
				v.logCancel()
			}
			return v, func() tea.Msg { return router.PopMsg{} }
		}

	case detailFetchDoneMsg:
		v.loading = false
		if msg.err != nil {
			return v, nil
		}
		overview := newOverviewTab().SetApp(msg.app)
		v.tabs = v.tabs.SetContent(0, overview)
		revisions := newRevisionsTab().SetRevisions(msg.revs)
		v.tabs = v.tabs.SetContent(1, revisions)
		settings := newSettingsTab().SetApp(msg.app)
		v.tabs = v.tabs.SetContent(3, settings)
		return v, nil

	case tabs.TabChangedMsg:
		if msg.Index == 2 && !v.logStarted {
			v.logStarted = true
			ctx, cancel := context.WithCancel(context.Background())
			v.logCancel = cancel
			p, rg, name := v.provider, v.rg, v.name
			return v, func() tea.Msg {
				ch, err := p.StreamLogs(ctx, rg, name)
				if err != nil {
					return logEntryMsg{done: true}
				}
				return logStreamReadyMsg{ch: ch}
			}
		}
		return v, nil

	case logStreamReadyMsg:
		lt := v.tabs.Tabs()[2].Content.(logsTab)
		lt, cmd := lt.StartStreaming(msg.ch)
		v.tabs = v.tabs.SetContent(2, lt)
		return v, cmd

	case logEntryMsg:
		tab := v.tabs.Tabs()[2].Content
		updated, cmd := tab.Update(msg)
		v.tabs = v.tabs.SetContent(2, updated)
		return v, cmd
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
		app, err := p.GetContainerApp(ctx, rg, name)
		if err != nil {
			return detailFetchDoneMsg{err: err}
		}
		revs, err := p.ListRevisions(ctx, rg, name)
		if err != nil {
			return detailFetchDoneMsg{app: app, err: err}
		}
		return detailFetchDoneMsg{app: app, revs: revs}
	}
}
