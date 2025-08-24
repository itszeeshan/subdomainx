package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/itszeeshan/subdomainx/internal/types"
)

// WriteCSV creates a CSV file with scan results
func WriteCSV(filename string, results *types.ScanResults) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Type",
		"Subdomain",
		"URL",
		"IP",
		"Port",
		"Protocol",
		"Status Code",
		"Title",
		"Technologies",
		"Content Length",
		"Source",
		"Service",
		"State",
		"Version",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write subdomain results
	for _, subdomain := range results.Subdomains {
		row := []string{
			"Subdomain",
			subdomain.Subdomain,
			"", // URL
			joinStrings(subdomain.IPs),
			"", // Port
			"", // Protocol
			"", // Status Code
			"", // Title
			"", // Technologies
			"", // Content Length
			subdomain.Source,
			"", // Service
			"", // State
			"", // Version
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write subdomain row: %v", err)
		}
	}

	// Write HTTP results
	for _, http := range results.HTTP {
		host, port, protocol, _ := parseURL(http.URL)

		row := []string{
			"HTTP",
			host,
			http.URL,
			"", // IP
			port,
			protocol,
			strconv.Itoa(http.StatusCode),
			http.Title,
			joinStrings(http.Technologies),
			strconv.Itoa(http.ContentLength),
			"httpx",
			"http",
			"open",
			"", // Version
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write HTTP row: %v", err)
		}
	}

	// Write port results
	for _, portResult := range results.Ports {
		for _, port := range portResult.Ports {
			row := []string{
				"Port",
				portResult.Host,
				"", // URL
				portResult.IP,
				strconv.Itoa(port.Number),
				port.Protocol,
				"", // Status Code
				"", // Title
				"", // Technologies
				"", // Content Length
				"smap",
				port.Service,
				port.State,
				port.Version,
			}
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write port row: %v", err)
			}
		}
	}

	return nil
}

// joinStrings joins a slice of strings with commas
func joinStrings(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += ", " + strs[i]
	}
	return result
}

// parseURL parses URL into components (reuse from burp_formatter.go)
// This function is already defined in burp_formatter.go
