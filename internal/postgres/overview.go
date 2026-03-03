package postgres

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type overviewTab struct {
	server provider.PostgresServer
	ready  bool
	width  int
	height int
}

func newOverviewTab() overviewTab {
	return overviewTab{}
}

func (t overviewTab) SetServer(s provider.PostgresServer) overviewTab {
	t.server = s
	t.ready = true
	return t
}

func (t overviewTab) Init() tea.Cmd { return nil }

func (t overviewTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		t.width = msg.Width
		t.height = msg.Height
	}
	return t, nil
}

func (t overviewTab) View() string {
	if !t.ready {
		return styles.Placeholder.Render("Loading...")
	}

	s := t.server
	labelStyle := lipgloss.NewStyle().Foreground(styles.Subtext).Width(18)
	valStyle := lipgloss.NewStyle().Foreground(styles.Text)

	row := func(label, val string) string {
		return labelStyle.Render(label) + valStyle.Render(val)
	}

	var lines []string
	lines = append(lines, row("Name", s.Name))
	lines = append(lines, row("State", s.State))
	lines = append(lines, row("Region", s.Region))
	lines = append(lines, row("Resource Group", s.ResourceGroup))
	lines = append(lines, row("Version", s.Version))
	lines = append(lines, row("SKU", s.SKU))
	lines = append(lines, row("Tier", s.Tier))
	lines = append(lines, row("Storage", fmt.Sprintf("%d GB", s.StorageGB)))
	lines = append(lines, row("Backup Retention", fmt.Sprintf("%d days", s.BackupRetention)))
	lines = append(lines, "")

	if s.FQDN != "" {
		headerStyle := lipgloss.NewStyle().Foreground(styles.Mauve).Bold(true)
		lines = append(lines, headerStyle.Render("Connection"))
		lines = append(lines, row("Host", s.FQDN))
		lines = append(lines, row("Port", "5432"))
		connStr := fmt.Sprintf("host=%s port=5432 sslmode=require", s.FQDN)
		lines = append(lines, row("Connection String", connStr))
	}

	return strings.Join(lines, "\n")
}
