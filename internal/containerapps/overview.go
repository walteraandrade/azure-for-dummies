package containerapps

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type overviewTab struct {
	app    provider.ContainerApp
	width  int
	height int
	ready  bool
}

func newOverviewTab() overviewTab {
	return overviewTab{}
}

func (t overviewTab) SetApp(app provider.ContainerApp) overviewTab {
	t.app = app
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

	a := t.app

	row := func(label, val string) string {
		return styles.DetailLabel.Render(label) + styles.DetailValue.Render(val)
	}

	var lines []string
	lines = append(lines, row("Name", a.Name))
	lines = append(lines, row("State", a.ProvisioningState))
	lines = append(lines, row("Region", a.Region))
	lines = append(lines, row("Resource Group", a.ResourceGroup))

	if a.FQDN != "" {
		lines = append(lines, row("FQDN", a.FQDN))
	}

	ingress := "Internal"
	if a.IngressExternal {
		ingress = fmt.Sprintf("External (port %d)", a.IngressPort)
	}
	lines = append(lines, row("Ingress", ingress))
	lines = append(lines, row("Scale", fmt.Sprintf("%d – %d", a.ScaleMin, a.ScaleMax)))
	lines = append(lines, "")

	if len(a.Containers) > 0 {
		lines = append(lines, styles.SectionHeader.Render("Containers"))

		for _, c := range a.Containers {
			lines = append(lines, row("  "+c.Name, fmt.Sprintf("%s  cpu:%.2f  mem:%s", c.Image, c.CPU, c.Memory)))
		}
	}

	return strings.Join(lines, "\n")
}
