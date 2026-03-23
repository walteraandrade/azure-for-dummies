package storage

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type containersTab struct {
	containers []provider.BlobContainer
	ready      bool
	width      int
	height     int
}

func newContainersTab() containersTab {
	return containersTab{}
}

func (t containersTab) SetContainers(c []provider.BlobContainer) containersTab {
	t.containers = c
	t.ready = true
	return t
}

func (t containersTab) Init() tea.Cmd { return nil }

func (t containersTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		t.width = msg.Width
		t.height = msg.Height
	}
	return t, nil
}

func (t containersTab) View() string {
	if !t.ready {
		return styles.Placeholder.Render("Loading containers...")
	}
	if len(t.containers) == 0 {
		return styles.Placeholder.Render("No containers found.")
	}

	nameCol := lipgloss.NewStyle().Foreground(styles.Text).Width(30)
	col := lipgloss.NewStyle().Foreground(styles.Subtext).Width(16)

	var lines []string
	lines = append(lines, styles.SectionHeader.Render("Blob Containers"))
	lines = append(lines, nameCol.Render("Name")+col.Render("Access")+col.Render("Lease")+col.Render("Last Modified"))
	lines = append(lines, strings.Repeat("─", 78))

	for _, c := range t.containers {
		modified := ""
		if !c.LastModified.IsZero() {
			modified = c.LastModified.Format("2006-01-02")
		}
		lines = append(lines, nameCol.Render(c.Name)+col.Render(c.PublicAccess)+col.Render(c.LeaseStatus)+col.Render(modified))
	}

	return strings.Join(lines, "\n")
}
