package statusbar

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/smarthow/azure-for-dummies/internal/styles"
)

type Model struct {
	subscription string
	user         string
	lastRefresh  time.Time
	err          error
	width        int
}

func New() Model {
	return Model{}
}

func (m Model) WithSubscription(sub string) Model {
	m.subscription = sub
	return m
}

func (m Model) WithUser(user string) Model {
	m.user = user
	return m
}

func (m Model) WithLastRefresh(t time.Time) Model {
	m.lastRefresh = t
	return m
}

func (m Model) WithError(err error) Model {
	m.err = err
	return m
}

func (m Model) WithWidth(w int) Model {
	m.width = w
	return m
}

func (m Model) View() string {
	if m.err != nil {
		errStr := fmt.Sprintf(" ✗ %s", m.err.Error())
		return styles.StatusBarError.Width(m.width).Render(errStr)
	}

	sub := m.subscription
	if sub == "" {
		sub = "—"
	}
	user := m.user
	if user == "" {
		user = "—"
	}

	subPart := styles.StatusBarKey.Render("SUB") + styles.StatusBarValue.Render(sub)
	userPart := styles.StatusBarKey.Render("USER") + styles.StatusBarValue.Render(user)

	var refreshPart string
	if !m.lastRefresh.IsZero() {
		refreshPart = styles.StatusBarKey.Render("REFRESHED") + styles.StatusBarValue.Render(m.lastRefresh.Format("15:04:05"))
	}

	left := lipgloss.JoinHorizontal(lipgloss.Top, subPart, userPart)
	if refreshPart != "" {
		left = lipgloss.JoinHorizontal(lipgloss.Top, left, refreshPart)
	}

	return styles.StatusBar.Width(m.width).Render(left)
}
