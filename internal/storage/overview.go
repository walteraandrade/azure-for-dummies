package storage

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type overviewTab struct {
	account provider.StorageAccount
	ready   bool
	width   int
	height  int
}

func newOverviewTab() overviewTab {
	return overviewTab{}
}

func (t overviewTab) SetAccount(a provider.StorageAccount) overviewTab {
	t.account = a
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

	a := t.account

	row := func(label, val string) string {
		return styles.DetailLabel.Render(label) + styles.DetailValue.Render(val)
	}

	var lines []string
	lines = append(lines, row("Name", a.Name))
	lines = append(lines, row("State", a.ProvisioningState))
	lines = append(lines, row("Region", a.Region))
	lines = append(lines, row("Resource Group", a.ResourceGroup))
	lines = append(lines, row("Kind", a.Kind))
	lines = append(lines, row("SKU", a.SKU))
	lines = append(lines, row("Access Tier", a.AccessTier))
	lines = append(lines, row("Min TLS", a.MinTLSVersion))
	lines = append(lines, row("HNS Enabled", fmt.Sprintf("%v", a.IsHnsEnabled)))
	lines = append(lines, row("Public Access", fmt.Sprintf("%v", a.AllowBlobPublicAccess)))
	lines = append(lines, row("Network Default", a.NetworkDefaultAction))
	lines = append(lines, "")

	endpoints := []struct{ label, val string }{
		{"Blob", a.PrimaryBlobEndpoint},
		{"File", a.PrimaryFileEndpoint},
		{"Table", a.PrimaryTableEndpoint},
		{"Queue", a.PrimaryQueueEndpoint},
	}

	hasEndpoint := false
	for _, e := range endpoints {
		if e.val != "" {
			hasEndpoint = true
			break
		}
	}
	if hasEndpoint {
		lines = append(lines, styles.SectionHeader.Render("Endpoints"))
		for _, e := range endpoints {
			if e.val != "" {
				lines = append(lines, row(e.label, e.val))
			}
		}
	}

	return strings.Join(lines, "\n")
}
