package notify

import (
	"fmt"
	"strings"
)

const maxListItems = 10

// FormatMarkdown formats a scan summary as Markdown (for Slack/Discord).
func FormatMarkdown(s ScanSummary) string {
	var b strings.Builder

	if s.Error != "" {
		fmt.Fprintf(&b, "**SubdomainX: Scan FAILED** for `%s`\n", s.Domain)
		fmt.Fprintf(&b, "Error: %s\n", s.Error)
		return b.String()
	}

	fmt.Fprintf(&b, "**SubdomainX: Scan completed** for `%s`\n", s.Domain)
	fmt.Fprintf(&b, "Scan ID: `%s` | Duration: %s\n", s.ScanID, s.Duration.Round(1e9))
	fmt.Fprintf(&b, "Subdomains: **%d**", s.TotalSubdomains)
	if s.TotalHTTP > 0 {
		fmt.Fprintf(&b, " | HTTP: **%d**", s.TotalHTTP)
	}
	if s.TotalPorts > 0 {
		fmt.Fprintf(&b, " | Ports: **%d**", s.TotalPorts)
	}
	b.WriteString("\n")

	if s.Diff != nil {
		b.WriteString("\n")
		formatDiffMarkdown(&b, s)
	}

	return b.String()
}

// FormatPlainText formats a scan summary as plain text (for Email/Telegram).
func FormatPlainText(s ScanSummary) string {
	var b strings.Builder

	if s.Error != "" {
		fmt.Fprintf(&b, "SubdomainX: Scan FAILED for %s\n", s.Domain)
		fmt.Fprintf(&b, "Error: %s\n", s.Error)
		return b.String()
	}

	fmt.Fprintf(&b, "SubdomainX: Scan completed for %s\n", s.Domain)
	fmt.Fprintf(&b, "Scan ID: %s | Duration: %s\n", s.ScanID, s.Duration.Round(1e9))
	fmt.Fprintf(&b, "Subdomains: %d", s.TotalSubdomains)
	if s.TotalHTTP > 0 {
		fmt.Fprintf(&b, " | HTTP: %d", s.TotalHTTP)
	}
	if s.TotalPorts > 0 {
		fmt.Fprintf(&b, " | Ports: %d", s.TotalPorts)
	}
	b.WriteString("\n")

	if s.Diff != nil {
		b.WriteString("\n")
		formatDiffPlainText(&b, s)
	}

	return b.String()
}

func formatDiffMarkdown(b *strings.Builder, s ScanSummary) {
	d := s.Diff
	fmt.Fprintf(b, "**Diff:** +%d new, -%d removed, ~%d changed\n", len(d.Added), len(d.Removed), len(d.IPChanges))

	if len(d.Added) > 0 {
		b.WriteString("\nNew subdomains:\n")
		writeList(b, d.Added, "+ `%s`\n")
	}
	if len(d.Removed) > 0 {
		b.WriteString("\nRemoved subdomains:\n")
		writeList(b, d.Removed, "- `%s`\n")
	}
	if len(d.IPChanges) > 0 {
		b.WriteString("\nIP changes:\n")
		for i, c := range d.IPChanges {
			if i >= maxListItems {
				fmt.Fprintf(b, "  ...and %d more\n", len(d.IPChanges)-maxListItems)
				break
			}
			fmt.Fprintf(b, "~ `%s` (%s -> %s)\n", c.Subdomain, strings.Join(c.OldIPs, ","), strings.Join(c.NewIPs, ","))
		}
	}
}

func formatDiffPlainText(b *strings.Builder, s ScanSummary) {
	d := s.Diff
	fmt.Fprintf(b, "Diff: +%d new, -%d removed, ~%d changed\n", len(d.Added), len(d.Removed), len(d.IPChanges))

	if len(d.Added) > 0 {
		b.WriteString("\nNew subdomains:\n")
		writeList(b, d.Added, "  + %s\n")
	}
	if len(d.Removed) > 0 {
		b.WriteString("\nRemoved subdomains:\n")
		writeList(b, d.Removed, "  - %s\n")
	}
	if len(d.IPChanges) > 0 {
		b.WriteString("\nIP changes:\n")
		for i, c := range d.IPChanges {
			if i >= maxListItems {
				fmt.Fprintf(b, "  ...and %d more\n", len(d.IPChanges)-maxListItems)
				break
			}
			fmt.Fprintf(b, "  ~ %s (%s -> %s)\n", c.Subdomain, strings.Join(c.OldIPs, ","), strings.Join(c.NewIPs, ","))
		}
	}
}

func writeList(b *strings.Builder, items []string, format string) {
	for i, item := range items {
		if i >= maxListItems {
			fmt.Fprintf(b, "  ...and %d more\n", len(items)-maxListItems)
			break
		}
		fmt.Fprintf(b, format, item)
	}
}
