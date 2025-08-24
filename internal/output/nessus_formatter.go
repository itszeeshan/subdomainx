package output

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/itszeeshan/subdomainx/internal/types"
)

// NessusReport represents the complete Nessus XML report
type NessusReport struct {
	XMLName xml.Name         `xml:"NessusClientData_v2"`
	Policy  NessusPolicy     `xml:"Policy"`
	Report  NessusReportData `xml:"Report"`
}

// NessusPolicy represents policy information in Nessus format
type NessusPolicy struct {
	XMLName     xml.Name `xml:"Policy"`
	PolicyName  string   `xml:"policyName"`
	PolicyID    string   `xml:"policyID"`
	PolicyUUID  string   `xml:"policyUUID"`
	PolicyOwner string   `xml:"policyOwner"`
	PolicyType  string   `xml:"policyType"`
}

// NessusReportData represents the report data in Nessus format
type NessusReportData struct {
	XMLName xml.Name     `xml:"Report"`
	Name    string       `xml:"name,attr"`
	Hosts   []NessusHost `xml:"ReportHost"`
}

// NessusHost represents a host in Nessus format
type NessusHost struct {
	XMLName xml.Name     `xml:"ReportHost"`
	Name    string       `xml:"name,attr"`
	Items   []NessusItem `xml:"ReportItem"`
}

// NessusItem represents an item in Nessus format
type NessusItem struct {
	XMLName             xml.Name `xml:"ReportItem"`
	Port                string   `xml:"port,attr"`
	SvcName             string   `xml:"svc_name,attr"`
	Protocol            string   `xml:"protocol,attr"`
	Severity            string   `xml:"severity,attr"`
	PluginID            string   `xml:"pluginID,attr"`
	PluginName          string   `xml:"pluginName,attr"`
	PluginFamily        string   `xml:"pluginFamily,attr"`
	PluginType          string   `xml:"plugin_type,attr"`
	PluginVersion       string   `xml:"plugin_version,attr"`
	RiskFactor          string   `xml:"risk_factor,attr"`
	Synopsis            string   `xml:"synopsis"`
	Description         string   `xml:"description"`
	Solution            string   `xml:"solution"`
	SeeAlso             string   `xml:"see_also"`
	CVE                 string   `xml:"cve"`
	BID                 string   `xml:"bid"`
	XRef                string   `xml:"xref"`
	PluginModDate       string   `xml:"plugin_mod_date"`
	PluginPubDate       string   `xml:"plugin_publication_date"`
	VulnPubDate         string   `xml:"vuln_publication_date"`
	PatchPubDate        string   `xml:"patch_publication_date"`
	CVSSVector          string   `xml:"cvss_vector"`
	CVSSBaseScore       string   `xml:"cvss_base_score"`
	CVSSTemporalScore   string   `xml:"cvss_temporal_score"`
	CVSSTemporalVector  string   `xml:"cvss_temporal_vector"`
	CVSS3Vector         string   `xml:"cvss3_vector"`
	CVSS3BaseScore      string   `xml:"cvss3_base_score"`
	CVSS3TemporalScore  string   `xml:"cvss3_temporal_score"`
	CVSS3TemporalVector string   `xml:"cvss3_temporal_vector"`
	MetasploitName      string   `xml:"metasploit_name"`
	ExploitAvailable    string   `xml:"exploit_available"`
	ExploitEase         string   `xml:"exploit_ease"`
	ExploitFrameworks   []string `xml:"exploit_frameworks>exploit_framework"`
	InTheNews           string   `xml:"in_the_news"`
	PluginOutput        string   `xml:"plugin_output"`
}

// WriteNessus creates a Nessus-compatible XML file
func WriteNessus(filename string, results *types.ScanResults) error {
	// Group results by host
	hostGroups := make(map[string][]types.HTTPResult)
	for _, http := range results.HTTP {
		host := extractHost(http.URL)
		hostGroups[host] = append(hostGroups[host], http)
	}

	// Create Nessus hosts
	var nessusHosts []NessusHost
	for host, httpResults := range hostGroups {
		var items []NessusItem

		for _, http := range httpResults {
			port, protocol := extractPortAndProtocol(http.URL)

			item := NessusItem{
				Port:          port,
				SvcName:       "http",
				Protocol:      protocol,
				Severity:      getSeverity(http.StatusCode),
				PluginID:      "99999", // Custom plugin ID for SubdomainX
				PluginName:    "SubdomainX - Discovered Web Service",
				PluginFamily:  "SubdomainX",
				PluginType:    "remote",
				PluginVersion: "1.0",
				RiskFactor:    getRiskFactor(http.StatusCode),
				Synopsis:      fmt.Sprintf("Web service discovered at %s", http.URL),
				Description:   generateDescription(http),
				Solution:      "Review discovered web services for security implications",
				SeeAlso:       "https://github.com/itszeeshan/subdomainx",
				PluginOutput:  generatePluginOutput(http),
			}
			items = append(items, item)
		}

		host := NessusHost{
			Name:  host,
			Items: items,
		}
		nessusHosts = append(nessusHosts, host)
	}

	// Create Nessus report
	report := NessusReport{
		Policy: NessusPolicy{
			PolicyName:  "SubdomainX Reconnaissance",
			PolicyID:    "subdomainx-001",
			PolicyUUID:  "subdomainx-policy-uuid",
			PolicyOwner: "SubdomainX",
			PolicyType:  "Custom",
		},
		Report: NessusReportData{
			Name:  "SubdomainX Scan Results",
			Hosts: nessusHosts,
		},
	}

	// Write XML file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create Nessus file: %v", err)
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("failed to encode Nessus XML: %v", err)
	}

	return nil
}

// extractHost extracts host from URL (reuse from zap_formatter.go)
// This function is already defined in zap_formatter.go

// extractPortAndProtocol extracts port and protocol from URL
func extractPortAndProtocol(url string) (string, string) {
	if len(url) > 8 && url[:8] == "https://" {
		return "443", "tcp"
	}
	return "80", "tcp"
}

// getSeverity determines severity based on status code
func getSeverity(statusCode int) string {
	if statusCode >= 500 {
		return "High"
	} else if statusCode >= 400 {
		return "Medium"
	} else if statusCode >= 300 {
		return "Low"
	}
	return "Info"
}

// getRiskFactor determines risk factor based on status code
func getRiskFactor(statusCode int) string {
	if statusCode >= 500 {
		return "High"
	} else if statusCode >= 400 {
		return "Medium"
	} else if statusCode >= 300 {
		return "Low"
	}
	return "None"
}

// generateDescription creates description for Nessus item
func generateDescription(http types.HTTPResult) string {
	desc := fmt.Sprintf("A web service was discovered at %s with status code %d.", http.URL, http.StatusCode)

	if http.Title != "" {
		desc += fmt.Sprintf(" The page title is: %s", http.Title)
	}

	if len(http.Technologies) > 0 {
		desc += fmt.Sprintf(" Technologies detected: %v", http.Technologies)
	}

	if http.ContentLength > 0 {
		desc += fmt.Sprintf(" Content length: %d bytes", http.ContentLength)
	}

	return desc
}

// generatePluginOutput creates plugin output for Nessus item
func generatePluginOutput(http types.HTTPResult) string {
	output := fmt.Sprintf("URL: %s\n", http.URL)
	output += fmt.Sprintf("Status Code: %d\n", http.StatusCode)

	if http.Title != "" {
		output += fmt.Sprintf("Title: %s\n", http.Title)
	}

	if len(http.Technologies) > 0 {
		output += fmt.Sprintf("Technologies: %v\n", http.Technologies)
	}

	if http.ContentLength > 0 {
		output += fmt.Sprintf("Content Length: %d\n", http.ContentLength)
	}

	return output
}
