package containerapps

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type settingsTab struct {
	app    provider.ContainerApp
	ready  bool
	width  int
	height int
}

func newSettingsTab() settingsTab {
	return settingsTab{}
}

func (t settingsTab) SetApp(app provider.ContainerApp) settingsTab {
	t.app = app
	t.ready = true
	return t
}

func (t settingsTab) Init() tea.Cmd { return nil }

func (t settingsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		t.width = msg.Width
		t.height = msg.Height
	}
	return t, nil
}

func (t settingsTab) View() string {
	if !t.ready {
		return styles.Placeholder.Render("Loading...")
	}

	var lines []string
	lines = append(lines, styles.SectionHeader.Render("Environment Variables"))

	for _, c := range t.app.Containers {
		if len(c.Env) == 0 {
			continue
		}
		lines = append(lines, "")
		lines = append(lines, lipgloss.NewStyle().Foreground(styles.Blue).Render(c.Name))
		for _, e := range c.Env {
			val := e.Value
			source := ""
			if e.SecretRef != "" {
				val = "••••••"
				source = " (secret: " + e.SecretRef + ")"
			}
			lines = append(lines, styles.DetailLabelWide.Render("  "+e.Name)+styles.DetailValue.Render(val)+styles.SecretValue.Render(source))
		}
	}

	if len(lines) == 1 {
		lines = append(lines, styles.Placeholder.Render("  No environment variables configured."))
	}

	return strings.Join(lines, "\n")
}
