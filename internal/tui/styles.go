package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorPrimary   = lipgloss.Color("#7C3AED") // violet
	colorSecondary = lipgloss.Color("#06B6D4") // cyan
	colorSuccess   = lipgloss.Color("#22C55E") // green
	colorWarning   = lipgloss.Color("#F59E0B") // amber
	colorError     = lipgloss.Color("#EF4444") // red
	colorMuted     = lipgloss.Color("#6B7280") // gray
	colorBg        = lipgloss.Color("#1E1E2E") // dark bg
	colorBorder    = lipgloss.Color("#44475A") // border gray

	// Tab styles
	activeTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(colorPrimary).
			Padding(0, 2)

	inactiveTab = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(0, 2)

	// Panel border
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	// Title inside panels
	panelTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSecondary)

	// Status indicators
	statusRunning   = lipgloss.NewStyle().Foreground(colorWarning)
	statusCompleted = lipgloss.NewStyle().Foreground(colorSuccess)
	statusFailed    = lipgloss.NewStyle().Foreground(colorError)
	statusSkipped   = lipgloss.NewStyle().Foreground(colorMuted)

	// Stats
	statLabel = lipgloss.NewStyle().Foreground(colorMuted)
	statValue = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))

	// Log levels
	logInfo  = lipgloss.NewStyle().Foreground(colorSecondary)
	logWarn  = lipgloss.NewStyle().Foreground(colorWarning)
	logError = lipgloss.NewStyle().Foreground(colorError)

	// Status bar
	statusBar = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(0, 1)

	// Spinner
	spinnerStyle = lipgloss.NewStyle().Foreground(colorPrimary)

	// Header
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Padding(0, 1)

	// Result table header
	tableHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSecondary).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder)

	// Stage indicator
	stageStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWarning)
)
