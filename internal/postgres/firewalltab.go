package postgres

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type firewallTab struct {
	rules  []provider.FirewallRuleInfo
	ready  bool
	width  int
	height int
}

func newFirewallTab() firewallTab {
	return firewallTab{}
}

func (t firewallTab) SetRules(rules []provider.FirewallRuleInfo) firewallTab {
	t.rules = rules
	t.ready = true
	return t
}

func (t firewallTab) Init() tea.Cmd { return nil }

func (t firewallTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		t.width = msg.Width
		t.height = msg.Height
	}
	return t, nil
}

func (t firewallTab) View() string {
	if !t.ready {
		return styles.Placeholder.Render("Loading firewall rules...")
	}
	if len(t.rules) == 0 {
		return styles.Placeholder.Render("No firewall rules configured.")
	}

	nameCol := lipgloss.NewStyle().Foreground(styles.Text).Width(30)
	ipCol := lipgloss.NewStyle().Foreground(styles.Subtext).Width(20)

	var lines []string
	lines = append(lines, styles.SectionHeader.Render("Firewall Rules"))
	lines = append(lines, nameCol.Render("Name")+ipCol.Render("Start IP")+ipCol.Render("End IP"))
	lines = append(lines, strings.Repeat("─", 70))

	for _, r := range t.rules {
		lines = append(lines, nameCol.Render(r.Name)+ipCol.Render(r.StartIP)+ipCol.Render(r.EndIP))
	}

	return strings.Join(lines, "\n")
}
