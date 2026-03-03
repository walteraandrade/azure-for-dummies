package containerapps

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type revisionItem struct {
	rev provider.RevisionInfo
}

func (i revisionItem) Title() string { return i.rev.Name }
func (i revisionItem) Description() string {
	active := "inactive"
	if i.rev.Active {
		active = "active"
	}
	return fmt.Sprintf("%s · traffic:%d%% · replicas:%d · %s · %s",
		active, i.rev.TrafficWeight, i.rev.Replicas,
		i.rev.HealthState, i.rev.CreatedTime.Format("2006-01-02 15:04"))
}
func (i revisionItem) FilterValue() string { return i.rev.Name }

type revisionsTab struct {
	list  list.Model
	ready bool
}

func newRevisionsTab() revisionsTab {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Revisions"
	l.Styles.Title = styles.ListTitle
	l.SetShowHelp(false)
	return revisionsTab{list: l}
}

func (t revisionsTab) SetRevisions(revs []provider.RevisionInfo) revisionsTab {
	items := make([]list.Item, len(revs))
	for i, r := range revs {
		items[i] = revisionItem{rev: r}
	}
	t.list.SetItems(items)
	t.ready = true
	return t
}

func (t revisionsTab) Init() tea.Cmd { return nil }

func (t revisionsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		t.list.SetSize(msg.Width, msg.Height)
		return t, nil
	}
	var cmd tea.Cmd
	t.list, cmd = t.list.Update(msg)
	return t, cmd
}

func (t revisionsTab) View() string {
	if !t.ready {
		return styles.Placeholder.Render("Loading revisions...")
	}
	return t.list.View()
}
