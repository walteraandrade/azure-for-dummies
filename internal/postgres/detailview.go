package postgres

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/azutil"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/router"
	"github.com/smarthow/azure-for-dummies/internal/styles"
	"github.com/smarthow/azure-for-dummies/internal/tabs"
)

type detailFetchDoneMsg struct {
	server provider.PostgresServer
	rules  []provider.FirewallRuleInfo
	series []provider.MetricSeries
	err    error
}

type detailView struct {
	id       string
	rg       string
	name     string
	provider provider.PostgresProvider
	tabs     tabs.Model
	spinner  spinner.Model
	loading  bool
	err      error
	width    int
	height   int
}

func newDetailView(id string, p provider.PostgresProvider) detailView {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = s.Style.Foreground(styles.Mauve)

	t := tabs.New(
		tabs.Tab{Title: "Overview", Content: newOverviewTab()},
		tabs.Tab{Title: "Firewall", Content: newFirewallTab()},
		tabs.Tab{Title: "Metrics", Content: newMetricsTab()},
	)

	return detailView{
		id:       id,
		rg:       azutil.ExtractRG(id),
		name:     azutil.ExtractName(id),
		provider: p,
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
			v.err = msg.err
			return v, nil
		}
		overview := newOverviewTab().SetServer(msg.server)
		v.tabs = v.tabs.SetContent(0, overview)
		fw := newFirewallTab().SetRules(msg.rules)
		v.tabs = v.tabs.SetContent(1, fw)
		metrics := newMetricsTab().SetMetrics(msg.series)
		v.tabs = v.tabs.SetContent(2, metrics)
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
	if v.err != nil {
		return styles.ErrorText.Render("Error: " + v.err.Error())
	}
	if v.loading {
		return v.spinner.View() + " Loading " + v.name + "..."
	}
	return v.tabs.View()
}

func (v detailView) fetchDetail() tea.Cmd {
	rg, name, id := v.rg, v.name, v.id
	p := v.provider
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server, err := p.GetServer(ctx, rg, name)
		if err != nil {
			return detailFetchDoneMsg{err: err}
		}
		rules, err := p.ListFirewallRules(ctx, rg, name)
		if err != nil {
			return detailFetchDoneMsg{server: server, err: err}
		}
		metricNames := []string{"cpu_percent", "memory_percent", "active_connections"}
		series, err := p.GetMetrics(ctx, id, metricNames)
		if err != nil {
			return detailFetchDoneMsg{server: server, rules: rules, err: err}
		}
		return detailFetchDoneMsg{server: server, rules: rules, series: series}
	}
}
