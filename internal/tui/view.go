package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.quitting {
		return "Shutting down...\n"
	}
	if m.width == 0 {
		return "Initializing..."
	}

	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	// Tabs
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	// Tab content
	switch m.activeTab {
	case tabDashboard:
		b.WriteString(m.renderDashboard())
	case tabResults:
		b.WriteString(m.renderResults())
	case tabLogs:
		b.WriteString(m.renderLogs())
	}

	// Status bar
	b.WriteString("\n")
	b.WriteString(m.renderStatusBar())

	return b.String()
}

func (m model) renderHeader() string {
	title := headerStyle.Render("SubdomainX Dashboard")
	elapsed := time.Since(m.startTime).Round(time.Second)
	right := statLabel.Render(fmt.Sprintf("Elapsed: %s", elapsed))

	gap := m.width - lipgloss.Width(title) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}
	return title + strings.Repeat(" ", gap) + right
}

func (m model) renderTabs() string {
	tabs := []string{"Dashboard", "Results", "Logs"}
	var rendered []string
	for i, t := range tabs {
		label := fmt.Sprintf(" %d %s ", i+1, t)
		if i == m.activeTab {
			rendered = append(rendered, activeTab.Render(label))
		} else {
			rendered = append(rendered, inactiveTab.Render(label))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func (m model) renderDashboard() string {
	// Left: Tool Progress, Right: Stats
	leftWidth := m.width*2/3 - 4
	rightWidth := m.width/3 - 4
	if leftWidth < 30 {
		leftWidth = 30
	}
	if rightWidth < 20 {
		rightWidth = 20
	}

	left := m.renderToolProgress(leftWidth)
	right := m.renderStats(rightWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right)
}

func (m model) renderToolProgress(width int) string {
	var lines []string
	lines = append(lines, panelTitle.Render("Tool Progress"))
	lines = append(lines, "")

	if len(m.tools) == 0 {
		lines = append(lines, statLabel.Render("  Waiting for tools to start..."))
	}

	for _, t := range m.tools {
		icon := m.toolIcon(t.status)
		name := fmt.Sprintf("%-15s", t.name)
		var detail string
		switch t.status {
		case "completed":
			detail = statusCompleted.Render(fmt.Sprintf("found %d", t.found))
		case "failed":
			detail = statusFailed.Render(truncate(t.err, 40))
		case "skipped":
			detail = statusSkipped.Render("skipped")
		default:
			detail = statusRunning.Render(fmt.Sprintf("scanning %s", t.domain))
		}
		lines = append(lines, fmt.Sprintf("  %s %s %s", icon, name, detail))
	}

	content := strings.Join(lines, "\n")
	return panelStyle.Width(width).Render(content)
}

func (m model) renderStats(width int) string {
	var lines []string
	lines = append(lines, panelTitle.Render("Stats"))
	lines = append(lines, "")

	lines = append(lines, m.statLine("Subdomains", m.totalSubdomains))
	lines = append(lines, m.statLine("HTTP Alive", m.totalHTTP))
	lines = append(lines, m.statLine("Ports", m.totalPorts))
	lines = append(lines, m.statLine("Takeover", m.totalTakeover))
	lines = append(lines, "")

	// Stage
	if m.scanDone {
		if m.scanError != nil {
			lines = append(lines, statusFailed.Render("SCAN FAILED"))
		} else {
			lines = append(lines, statusCompleted.Render("SCAN COMPLETE"))
		}
	} else {
		lines = append(lines, stageStyle.Render(fmt.Sprintf("Stage: %s", m.currentStage)))
		lines = append(lines, m.spinner.View()+" "+statLabel.Render(m.stageMessage))
	}

	lines = append(lines, "")

	// Resources
	lines = append(lines, panelTitle.Render("Resources"))
	lines = append(lines, fmt.Sprintf("  %s %s  %s %s",
		statLabel.Render("Mem:"),
		statValue.Render(fmt.Sprintf("%.1fMB", m.memoryMB)),
		statLabel.Render("Goroutines:"),
		statValue.Render(fmt.Sprintf("%d", m.goroutines)),
	))

	content := strings.Join(lines, "\n")
	return panelStyle.Width(width).Render(content)
}

func (m model) statLine(label string, value int) string {
	return fmt.Sprintf("  %s %s",
		statLabel.Render(fmt.Sprintf("%-12s", label+":")),
		statValue.Render(fmt.Sprintf("%d", value)),
	)
}

func (m model) renderResults() string {
	results := m.getFilteredResults()
	var lines []string

	// Header
	header := fmt.Sprintf("Results: %d subdomains", len(results))
	if m.filterActive {
		header += fmt.Sprintf("  Filter: %s_", m.filterText)
	} else if m.filterText != "" {
		header += fmt.Sprintf("  [filtered: %q]", m.filterText)
	}
	lines = append(lines, panelTitle.Render(header))
	lines = append(lines, "")

	// Column headers
	sortIndicator := []string{"  ", "  ", "  "}
	if m.sortAsc {
		sortIndicator[m.sortColumn] = " ^"
	} else {
		sortIndicator[m.sortColumn] = " v"
	}

	colHeader := fmt.Sprintf("  %-50s %-20s %-20s",
		"Subdomain"+sortIndicator[0],
		"Source"+sortIndicator[1],
		"IPs"+sortIndicator[2],
	)
	lines = append(lines, tableHeader.Render(colHeader))

	// Rows
	visible := m.visibleRows()
	end := m.resultOffset + visible
	if end > len(results) {
		end = len(results)
	}

	for i := m.resultOffset; i < end; i++ {
		r := results[i]
		ips := strings.Join(r.IPs, ", ")
		if ips == "" {
			ips = "-"
		}
		line := fmt.Sprintf("  %-50s %-20s %-20s",
			truncate(r.Subdomain, 48),
			truncate(r.Source, 18),
			truncate(ips, 18),
		)
		lines = append(lines, line)
	}

	if len(results) == 0 {
		lines = append(lines, statLabel.Render("  No results yet..."))
	}

	content := strings.Join(lines, "\n")
	return panelStyle.Width(m.width - 4).Render(content)
}

func (m model) renderLogs() string {
	header := panelTitle.Render(fmt.Sprintf("Logs (%d entries)", len(m.logs)))
	content := header + "\n\n" + m.logViewport.View()
	return panelStyle.Width(m.width - 4).Height(m.height - 8).Render(content)
}

func (m model) renderStatusBar() string {
	var parts []string

	parts = append(parts, "[Tab] Switch")
	if m.activeTab == tabResults {
		parts = append(parts, "[s] Sort", "[/] Filter", "[j/k] Scroll")
	}
	if m.activeTab == tabLogs {
		parts = append(parts, "[j/k] Scroll")
	}
	parts = append(parts, "[e] Export", "[q] Quit")

	return statusBar.Render(strings.Join(parts, "  "))
}

func (m model) toolIcon(status string) string {
	switch status {
	case "completed":
		return statusCompleted.Render("*")
	case "failed":
		return statusFailed.Render("x")
	case "skipped":
		return statusSkipped.Render("-")
	default:
		return m.spinner.View()
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "~"
}
