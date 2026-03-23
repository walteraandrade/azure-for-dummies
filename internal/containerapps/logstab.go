package containerapps

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/smarthow/azure-for-dummies/internal/provider"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type logEntryMsg struct {
	entry provider.LogEntry
	done  bool
}

type logsTab struct {
	viewport viewport.Model
	lines    []string
	ch       <-chan provider.LogEntry
	started  bool
	errMsg   string
	width    int
	height   int
}

func newLogsTab() logsTab {
	return logsTab{}
}

func (t logsTab) Init() tea.Cmd { return nil }

func (t logsTab) StartStreaming(ch <-chan provider.LogEntry) (logsTab, tea.Cmd) {
	t.ch = ch
	t.started = true
	t.viewport = viewport.New(t.width, t.height)
	return t, waitForLog(ch)
}

func (t logsTab) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height
		t.viewport.Width = msg.Width
		t.viewport.Height = msg.Height
		return t, nil

	case logEntryMsg:
		if msg.done {
			return t, nil
		}
		line := fmt.Sprintf("%s  %s", msg.entry.Timestamp.Format("15:04:05"), msg.entry.Message)
		t.lines = append(t.lines, line)
		t.viewport.SetContent(strings.Join(t.lines, "\n"))
		t.viewport.GotoBottom()
		return t, waitForLog(t.ch)
	}

	var cmd tea.Cmd
	t.viewport, cmd = t.viewport.Update(msg)
	return t, cmd
}

func (t logsTab) View() string {
	if t.errMsg != "" {
		return styles.Placeholder.Render("Log stream error: " + t.errMsg)
	}
	if !t.started {
		return styles.Placeholder.Render("Press tab to activate log streaming...")
	}
	if len(t.lines) == 0 {
		return styles.Placeholder.Render("Waiting for logs...")
	}
	return t.viewport.View()
}

func waitForLog(ch <-chan provider.LogEntry) tea.Cmd {
	if ch == nil {
		return nil
	}
	return func() tea.Msg {
		entry, ok := <-ch
		if !ok {
			return logEntryMsg{done: true}
		}
		return logEntryMsg{entry: entry}
	}
}
