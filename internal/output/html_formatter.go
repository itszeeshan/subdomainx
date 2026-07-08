package output

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/types"
)

//go:embed templates/report.html
var reportTemplateStr string

const itemsPerPage = 50

// statEntry is a name/count pair for source and domain breakdown charts.
type statEntry struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// reportData holds all values injected into the HTML report template.
type reportData struct {
	ScanID         string
	GeneratedAt    string
	WildcardFile   string
	Threads        int
	Retries        int
	Timeout        int
	RateLimit      int
	SubdomainCount int
	HTTPCount      int
	PortCount      int
	UniqueSources  int
	UniqueDomains  int
	Subdomains     template.JS // safe JSON array
	HTTPResults    template.JS // safe JSON array
	PortResults    template.JS // safe JSON array
	DomainStats    template.JS // [{name, count}]
	SourceStats    template.JS // [{name, count}]
	ItemsPerPage   int
}

// WriteHTML renders the embedded report template with the scan results and
// writes it to filename.
func WriteHTML(filename string, cfg *config.Config, results *types.ScanResults) error {
	tmpl, err := template.New("report").Parse(reportTemplateStr)
	if err != nil {
		return err
	}

	domainStats, sourcesStats := computeStats(results.Subdomains)

	data := reportData{
		ScanID:         cfg.UniqueName,
		GeneratedAt:    time.Now().Format("January 2, 2006 at 15:04:05"),
		WildcardFile:   cfg.WildcardFile,
		Threads:        cfg.Threads,
		Retries:        cfg.Retries,
		Timeout:        cfg.Timeout,
		RateLimit:      cfg.RateLimit,
		SubdomainCount: len(results.Subdomains),
		HTTPCount:      len(results.HTTP),
		PortCount:      len(results.Ports),
		UniqueSources:  len(sourcesStats),
		UniqueDomains:  len(domainStats),
		Subdomains:     buildSubdomainRows(results.Subdomains),
		HTTPResults:    buildHTTPRows(results.HTTP),
		PortResults:    buildPortRows(results.Ports),
		DomainStats:    marshalJS(domainStats),
		SourceStats:    marshalJS(sourcesStats),
		ItemsPerPage:   itemsPerPage,
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// computeStats returns domain and source breakdown sorted by count descending.
func computeStats(subdomains []types.SubdomainResult) ([]statEntry, []statEntry) {
	domainCounts := make(map[string]int)
	sourceCounts := make(map[string]int)

	for _, s := range subdomains {
		parent := extractParent(s.Subdomain)
		domainCounts[parent]++
		for _, src := range strings.Split(s.Source, ",") {
			src = strings.TrimSpace(src)
			if src != "" {
				sourceCounts[src]++
			}
		}
	}

	return toSortedEntries(domainCounts), toSortedEntries(sourceCounts)
}

func toSortedEntries(m map[string]int) []statEntry {
	entries := make([]statEntry, 0, len(m))
	for k, v := range m {
		entries = append(entries, statEntry{Name: k, Count: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Count > entries[j].Count
	})
	return entries
}

// extractParent returns the registrable domain portion (last two labels) of a
// fully-qualified subdomain name, e.g. "api.palantir.com" → "palantir.com".
func extractParent(subdomain string) string {
	parts := strings.Split(subdomain, ".")
	if len(parts) <= 2 {
		return subdomain
	}
	return strings.Join(parts[len(parts)-2:], ".")
}

// buildSubdomainRows marshals subdomain results to a JSON array safe for
// embedding inside a <script> block.
func buildSubdomainRows(subdomains []types.SubdomainResult) template.JS {
	type row struct {
		Subdomain string `json:"subdomain"`
		Parent    string `json:"parent"`
		Source    string `json:"source"`
		IPs       string `json:"ips"`
	}
	rows := make([]row, 0, len(subdomains))
	for _, s := range subdomains {
		ips := strings.Join(s.IPs, ", ")
		if ips == "" {
			ips = "N/A"
		}
		rows = append(rows, row{
			Subdomain: s.Subdomain,
			Parent:    extractParent(s.Subdomain),
			Source:    s.Source,
			IPs:       ips,
		})
	}
	return marshalJS(rows)
}

// buildHTTPRows marshals HTTP results to a JSON array safe for embedding
// inside a <script> block.
func buildHTTPRows(httpResults []types.HTTPResult) template.JS {
	type row struct {
		URL           string `json:"url"`
		Status        int    `json:"status"`
		StatusClass   string `json:"statusClass"`
		Title         string `json:"title"`
		ContentLength int    `json:"contentLength"`
		Technologies  string `json:"technologies"`
	}
	rows := make([]row, 0, len(httpResults))
	for _, r := range httpResults {
		tech := strings.Join(r.Technologies, ", ")
		if tech == "" {
			tech = "N/A"
		}
		rows = append(rows, row{
			URL:           r.URL,
			Status:        r.StatusCode,
			StatusClass:   statusClass(r.StatusCode),
			Title:         r.Title,
			ContentLength: r.ContentLength,
			Technologies:  tech,
		})
	}
	return marshalJS(rows)
}

// buildPortRows marshals port scan results to a JSON array safe for embedding
// inside a <script> block.
func buildPortRows(portResults []types.PortResult) template.JS {
	type row struct {
		Host     string `json:"host"`
		IP       string `json:"ip"`
		Port     int    `json:"port"`
		Protocol string `json:"protocol"`
		State    string `json:"state"`
		Service  string `json:"service"`
		Version  string `json:"version"`
	}
	var rows []row
	for _, pr := range portResults {
		for _, p := range pr.Ports {
			rows = append(rows, row{
				Host:     pr.Host,
				IP:       pr.IP,
				Port:     p.Number,
				Protocol: p.Protocol,
				State:    p.State,
				Service:  p.Service,
				Version:  p.Version,
			})
		}
	}
	return marshalJS(rows)
}

// marshalJS marshals v to JSON and returns it as template.JS, escaping any
// "</script>" sequence to prevent breaking out of the enclosing script block.
func marshalJS(v any) template.JS {
	b, _ := json.Marshal(v)
	safe := strings.ReplaceAll(string(b), "</script>", `<\/script>`)
	return template.JS(safe)
}

func statusClass(code int) string {
	return "status-" + strconv.Itoa(code)
}
