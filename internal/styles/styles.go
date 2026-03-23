package styles

import "github.com/charmbracelet/lipgloss"

// Catppuccin Mocha palette
var (
	Base    = lipgloss.Color("#1e1e2e")
	Mantle  = lipgloss.Color("#181825")
	Surface = lipgloss.Color("#313244")
	Overlay = lipgloss.Color("#6c7086")
	Text    = lipgloss.Color("#cdd6f4")
	Subtext = lipgloss.Color("#a6adc8")
	Mauve   = lipgloss.Color("#cba6f7")
	Blue    = lipgloss.Color("#89b4fa")
	Green   = lipgloss.Color("#a6e3a1")
	Red     = lipgloss.Color("#f38ba8")
	Yellow  = lipgloss.Color("#f9e2af")
)

var StatusBar = lipgloss.NewStyle().
	Background(Mantle).
	Foreground(Subtext).
	Padding(0, 1)

var StatusBarKey = lipgloss.NewStyle().
	Background(Surface).
	Foreground(Mauve).
	Bold(true).
	Padding(0, 1)

var StatusBarValue = lipgloss.NewStyle().
	Background(Mantle).
	Foreground(Text).
	Padding(0, 1)

var StatusBarError = lipgloss.NewStyle().
	Background(Mantle).
	Foreground(Red).
	Padding(0, 1)

var ContentArea = lipgloss.NewStyle().
	Foreground(Text)

var Placeholder = lipgloss.NewStyle().
	Foreground(Overlay).
	Italic(true)

var ModuleCard = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Surface).
	Foreground(Text).
	Padding(1, 2)

var ModuleCardSelected = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(Mauve).
	Foreground(Mauve).
	Bold(true).
	Padding(1, 2)

var ListTitle = lipgloss.NewStyle().
	Foreground(Mauve).
	Bold(true)

var TabActive = lipgloss.NewStyle().
	Foreground(Mauve).
	Bold(true).
	Padding(0, 2)

var TabInactive = lipgloss.NewStyle().
	Foreground(Subtext).
	Padding(0, 2)

var TabBar = lipgloss.NewStyle().
	BorderBottom(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(Surface)

var ErrorText = lipgloss.NewStyle().Foreground(Red)

var DetailLabel = lipgloss.NewStyle().Foreground(Subtext).Width(18)

var DetailLabelWide = lipgloss.NewStyle().Foreground(Subtext).Width(24)

var DetailValue = lipgloss.NewStyle().Foreground(Text)

var SectionHeader = lipgloss.NewStyle().Foreground(Mauve).Bold(true)

var SecretValue = lipgloss.NewStyle().Foreground(Overlay)
