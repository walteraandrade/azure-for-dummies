package containerapps

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/auth"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/router"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type appItem struct {
	app provider.ContainerApp
}

func (i appItem) Title() string { return i.app.Name }
func (i appItem) Description() string {
	return fmt.Sprintf("%s · %s · %s · %s", i.app.State, i.app.ResourceGroup, i.app.Region, i.app.Revision)
}
func (i appItem) FilterValue() string { return i.app.Name }

type listView struct {
	list     list.Model
	spinner  spinner.Model
	loading  bool
	provider *azureProvider
	auth     *auth.Context
}

func newListView(p *azureProvider, a *auth.Context) listView {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Container Apps"
	l.Styles.Title = styles.ListTitle

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = s.Style.Foreground(styles.Mauve)

	return listView{list: l, spinner: s, loading: true, provider: p, auth: a}
}

func (v listView) Init() tea.Cmd {
	return v.spinner.Tick
}

func (v listView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.list.SetSize(msg.Width, msg.Height)

	case FetchDoneMsg:
		v.loading = false
		items := make([]list.Item, len(msg.Apps))
		for i, a := range msg.Apps {
			items[i] = appItem{app: a}
		}
		cmd := v.list.SetItems(items)
		return v, cmd

	case tea.KeyMsg:
		if msg.String() == "esc" {
			if v.list.FilterState() == list.Filtering {
				break
			}
			return v, func() tea.Msg { return router.PopMsg{} }
		}
		if msg.String() == "enter" && !v.loading {
			if v.list.FilterState() != list.Filtering {
				if sel, ok := v.list.SelectedItem().(appItem); ok {
					return v, func() tea.Msg {
						return router.PushMsg{Screen: newDetailView(sel.app.ID, v.provider, v.auth)}
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
		return v.spinner.View() + " Loading container apps..."
	}
	return v.list.View()
}
