package postgres

import (
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

var sparkChars = []rune("▁▂▃▄▅▆▇█")

type metricsTab struct {
	series []provider.MetricSeries
	ready  bool
	width  int
	height int
}

func newMetricsTab() metricsTab {
	return metricsTab{}
}

func (t metricsTab) SetMetrics(series []provider.MetricSeries) metricsTab {
	t.series = series
	t.ready = true
	return t
}

func (t metricsTab) Init() tea.Cmd { return nil }

func (t metricsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		t.width = msg.Width
		t.height = msg.Height
	}
	return t, nil
}

func (t metricsTab) View() string {
	if !t.ready {
		return styles.Placeholder.Render("Loading metrics...")
	}
	if len(t.series) == 0 {
		return styles.Placeholder.Render("No metrics available.")
	}

	headerStyle := lipgloss.NewStyle().Foreground(styles.Mauve).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(styles.Subtext).Width(22)
	valStyle := lipgloss.NewStyle().Foreground(styles.Text)
	sparkStyle := lipgloss.NewStyle().Foreground(styles.Green)

	var lines []string
	lines = append(lines, headerStyle.Render("Metrics (last 1h, 5min intervals)"))
	lines = append(lines, "")

	for _, ms := range t.series {
		current := float64(0)
		if len(ms.Points) > 0 {
			current = ms.Points[len(ms.Points)-1].Average
		}
		spark := renderSparkline(ms.Points)
		lines = append(lines,
			labelStyle.Render(ms.Name)+
				sparkStyle.Render(spark)+"  "+
				valStyle.Render(fmt.Sprintf("%.1f %s", current, ms.Unit)),
		)
	}

	return strings.Join(lines, "\n")
}

func renderSparkline(points []provider.MetricPoint) string {
	if len(points) == 0 {
		return ""
	}

	minVal, maxVal := math.MaxFloat64, -math.MaxFloat64
	for _, p := range points {
		if p.Average < minVal {
			minVal = p.Average
		}
		if p.Average > maxVal {
			maxVal = p.Average
		}
	}

	span := maxVal - minVal
	if span == 0 {
		span = 1
	}

	var sb strings.Builder
	for _, p := range points {
		idx := int((p.Average - minVal) / span * float64(len(sparkChars)-1))
		if idx >= len(sparkChars) {
			idx = len(sparkChars) - 1
		}
		if idx < 0 {
			idx = 0
		}
		sb.WriteRune(sparkChars[idx])
	}
	return sb.String()
}
