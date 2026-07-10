package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/itszeeshan/subdomainx/internal/types"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.logViewport.Width = msg.Width - 4
		m.logViewport.Height = msg.Height - 8
		return m, nil

	case tea.KeyMsg:
		// When filter is active, handle text input
		if m.filterActive {
			switch msg.String() {
			case "esc":
				m.filterActive = false
				m.filterText = ""
			case "enter":
				m.filterActive = false
			case "backspace":
				if len(m.filterText) > 0 {
					m.filterText = m.filterText[:len(m.filterText)-1]
				}
			default:
				if len(msg.String()) == 1 {
					m.filterText += msg.String()
				}
			}
			return m, nil
		}

		switch {
		case msg.String() == "q" || msg.String() == "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case msg.String() == "tab":
			m.activeTab = (m.activeTab + 1) % tabCount
		case msg.String() == "shift+tab":
			m.activeTab = (m.activeTab - 1 + tabCount) % tabCount
		case msg.String() == "1":
			m.activeTab = tabDashboard
		case msg.String() == "2":
			m.activeTab = tabResults
		case msg.String() == "3":
			m.activeTab = tabLogs
		case msg.String() == "j" || msg.String() == "down":
			switch m.activeTab {
			case tabResults:
				m.resultOffset++
				maxOffset := len(m.getFilteredResults()) - m.visibleRows()
				if maxOffset < 0 {
					maxOffset = 0
				}
				if m.resultOffset > maxOffset {
					m.resultOffset = maxOffset
				}
			case tabLogs:
				m.logViewport.ScrollDown(1)
			}
		case msg.String() == "k" || msg.String() == "up":
			switch m.activeTab {
			case tabResults:
				m.resultOffset--
				if m.resultOffset < 0 {
					m.resultOffset = 0
				}
			case tabLogs:
				m.logViewport.ScrollUp(1)
			}
		case msg.String() == "s":
			if m.activeTab == tabResults {
				m.sortColumn = (m.sortColumn + 1) % 3
				m.sortResults()
			}
		case msg.String() == "/":
			if m.activeTab == tabResults {
				m.filterActive = true
				m.filterText = ""
			}
		case msg.String() == "e":
			// Export handled at the run.go level via checkpoint
			m.logs = append(m.logs, LogMsg{
				Level:   "info",
				Message: "Results are auto-saved to the output directory",
				Time:    time.Now(),
			})
			m.updateLogViewport()
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case TickMsg:
		cmds = append(cmds, tickCmd())

	// Pipeline events
	case StageMsg:
		m.currentStage = msg.Stage
		m.stageMessage = msg.Message
		m.logs = append(m.logs, LogMsg{
			Level:   "info",
			Message: fmt.Sprintf("[%s] %s", msg.Stage, msg.Message),
			Time:    time.Now(),
		})
		m.updateLogViewport()

	case ToolProgressMsg:
		m.updateToolStatus(msg)
		if msg.Status != "running" {
			m.logs = append(m.logs, LogMsg{
				Level:   toolLogLevel(msg.Status),
				Message: formatToolLog(msg),
				Time:    time.Now(),
			})
			m.updateLogViewport()
		}

	case ResultMsg:
		m.subdomains = msg.Results
		m.totalSubdomains = msg.Total

	case HTTPResultMsg:
		m.httpResults = msg.Results
		m.totalHTTP = msg.Total

	case PortResultMsg:
		m.portResults = msg.Results
		m.totalPorts = msg.Total

	case TakeoverResultMsg:
		m.takeoverResults = msg.Results
		m.totalTakeover = len(msg.Results)

	case LogMsg:
		m.logs = append(m.logs, msg)
		m.updateLogViewport()

	case ResourceMsg:
		m.memoryMB = msg.MemoryMB
		m.goroutines = msg.Goroutines

	case ScanCompleteMsg:
		m.scanDone = true
		m.scanError = msg.Error
		level := "info"
		message := "Scan completed successfully"
		if msg.Error != nil {
			level = "error"
			message = fmt.Sprintf("Scan failed: %v", msg.Error)
		}
		m.logs = append(m.logs, LogMsg{Level: level, Message: message, Time: time.Now()})
		m.updateLogViewport()
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateToolStatus(msg ToolProgressMsg) {
	for i, t := range m.tools {
		if t.name == msg.Tool && t.domain == msg.Domain {
			m.tools[i].status = msg.Status
			m.tools[i].found = msg.Found
			m.tools[i].err = msg.Error
			return
		}
	}
	m.tools = append(m.tools, toolStatus{
		name:   msg.Tool,
		domain: msg.Domain,
		status: msg.Status,
		found:  msg.Found,
		err:    msg.Error,
	})
}

func (m *model) sortResults() {
	sort.Slice(m.subdomains, func(i, j int) bool {
		var a, b string
		switch m.sortColumn {
		case 0:
			a, b = m.subdomains[i].Subdomain, m.subdomains[j].Subdomain
		case 1:
			a, b = m.subdomains[i].Source, m.subdomains[j].Source
		case 2:
			a = strings.Join(m.subdomains[i].IPs, ",")
			b = strings.Join(m.subdomains[j].IPs, ",")
		}
		if m.sortAsc {
			return a < b
		}
		return a > b
	})
}

func (m model) getFilteredResults() []types.SubdomainResult {
	if m.filterText == "" {
		return m.subdomains
	}
	filter := strings.ToLower(m.filterText)
	var filtered []types.SubdomainResult
	for _, r := range m.subdomains {
		if strings.Contains(strings.ToLower(r.Subdomain), filter) ||
			strings.Contains(strings.ToLower(r.Source), filter) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func (m model) visibleRows() int {
	rows := m.height - 10
	if rows < 1 {
		return 1
	}
	return rows
}

func (m *model) updateLogViewport() {
	var lines []string
	for _, l := range m.logs {
		ts := l.Time.Format("15:04:05")
		var levelStr string
		switch l.Level {
		case "error":
			levelStr = logError.Render("ERR")
		case "warn":
			levelStr = logWarn.Render("WRN")
		default:
			levelStr = logInfo.Render("INF")
		}
		lines = append(lines, fmt.Sprintf("%s %s %s", ts, levelStr, l.Message))
	}
	m.logViewport.SetContent(strings.Join(lines, "\n"))
	m.logViewport.GotoBottom()
}

func toolLogLevel(status string) string {
	switch status {
	case "failed":
		return "error"
	case "skipped":
		return "warn"
	default:
		return "info"
	}
}

func formatToolLog(msg ToolProgressMsg) string {
	switch msg.Status {
	case "completed":
		return fmt.Sprintf("%s found %d subdomains for %s", msg.Tool, msg.Found, msg.Domain)
	case "failed":
		return fmt.Sprintf("%s failed for %s: %s", msg.Tool, msg.Domain, msg.Error)
	case "skipped":
		return fmt.Sprintf("%s skipped (not in PATH)", msg.Tool)
	default:
		return fmt.Sprintf("%s %s for %s", msg.Tool, msg.Status, msg.Domain)
	}
}
