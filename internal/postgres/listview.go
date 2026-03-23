package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/router"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type serverItem struct {
	server provider.PostgresServer
}

func (i serverItem) Title() string { return i.server.Name }
func (i serverItem) Description() string {
	return fmt.Sprintf("%s · %s · %s · %s", i.server.Tier, i.server.State, i.server.Region, i.server.Version)
}
func (i serverItem) FilterValue() string { return i.server.Name }

type listView struct {
	list     list.Model
	spinner  spinner.Model
	loading  bool
	err      error
	provider provider.PostgresProvider
}

func newListView(p provider.PostgresProvider) listView {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "PostgreSQL Flexible Servers"
	l.Styles.Title = styles.ListTitle

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = s.Style.Foreground(styles.Mauve)

	return listView{list: l, spinner: s, loading: true, provider: p}
}

func (v listView) Init() tea.Cmd {
	return tea.Batch(v.spinner.Tick, v.fetchList())
}

func (v listView) fetchList() tea.Cmd {
	p := v.provider
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		servers, err := p.ListServers(ctx)
		return FetchDoneMsg{Servers: servers, Err: err}
	}
}

func (v listView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.list.SetSize(msg.Width, msg.Height)

	case FetchDoneMsg:
		v.loading = false
		if msg.Err != nil {
			v.err = msg.Err
			return v, nil
		}
		items := make([]list.Item, len(msg.Servers))
		for i, s := range msg.Servers {
			items[i] = serverItem{server: s}
		}
		cmd := v.list.SetItems(items)
		return v, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if v.list.FilterState() != list.Filtering {
				return v, tea.Quit
			}
		case "esc":
			if v.list.FilterState() == list.Filtering {
				break
			}
			return v, func() tea.Msg { return router.PopMsg{} }
		case "enter":
			if !v.loading && v.list.FilterState() != list.Filtering {
				if sel, ok := v.list.SelectedItem().(serverItem); ok {
					return v, func() tea.Msg {
						return router.PushMsg{Screen: newDetailView(sel.server.ID, v.provider)}
					}
				}
			}
		}
	}

	if v.loading {
		var cmd tea.Cmd
		v.spinner, cmd = v.spinner.Update(msg)
		return v, cmd
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v, cmd
}

func (v listView) View() string {
	if v.loading {
		return v.spinner.View() + " Loading PostgreSQL servers..."
	}
	if v.err != nil {
		return styles.ErrorText.Render("Error: " + v.err.Error())
	}
	return v.list.View()
}
