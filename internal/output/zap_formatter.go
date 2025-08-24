package output

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/itszeeshan/subdomainx/internal/types"
)

// ZAPSite represents a site in ZAP XML format
type ZAPSite struct {
	XMLName xml.Name `xml:"site"`
	Name    string   `xml:"name,attr"`
	Host    string   `xml:"host,attr"`
	Port    string   `xml:"port,attr"`
	SSL     string   `xml:"ssl,attr"`
	URLs    []ZAPURL `xml:"urls>url"`
}

// ZAPURL represents a URL in ZAP XML format
type ZAPURL struct {
	XMLName xml.Name `xml:"url"`
	Method  string   `xml:"method,attr"`
	URL     string   `xml:",chardata"`
}

// ZAPReport represents the complete ZAP XML report
type ZAPReport struct {
	XMLName   xml.Name    `xml:"zap-report"`
	Generated string      `xml:"generated,attr"`
	Version   string      `xml:"version,attr"`
	ScanInfo  ZAPScanInfo `xml:"scan-info"`
	Sites     []ZAPSite   `xml:"sites>site"`
}

// ZAPScanInfo represents scan information in ZAP format
type ZAPScanInfo struct {
	XMLName   xml.Name `xml:"scan-info"`
	StartTime string   `xml:"start-time,attr"`
	EndTime   string   `xml:"end-time,attr"`
	ScanType  string   `xml:"scan-type,attr"`
	TargetURL string   `xml:"target-url,attr"`
}

// WriteZAP creates a ZAP-compatible XML file
func WriteZAP(filename string, results *types.ScanResults) error {
	// Group HTTP results by host
	hostGroups := make(map[string][]types.HTTPResult)
	for _, http := range results.HTTP {
		host := extractHost(http.URL)
		hostGroups[host] = append(hostGroups[host], http)
	}

	// Create ZAP sites
	var zapSites []ZAPSite
	for host, httpResults := range hostGroups {
		var urls []ZAPURL
		for _, http := range httpResults {
			url := ZAPURL{
				Method: "GET",
				URL:    http.URL,
			}
			urls = append(urls, url)
		}

		// Determine port and SSL
		port := "80"
		ssl := "false"
		if len(httpResults) > 0 {
			port, ssl = extractPortAndSSL(httpResults[0].URL)
		}

		site := ZAPSite{
			Name: host,
			Host: host,
			Port: port,
			SSL:  ssl,
			URLs: urls,
		}
		zapSites = append(zapSites, site)
	}

	// Create ZAP report
	now := time.Now()
	report := ZAPReport{
		Generated: now.Format(time.RFC3339),
		Version:   "2.0",
		ScanInfo: ZAPScanInfo{
			StartTime: now.Format(time.RFC3339),
			EndTime:   now.Format(time.RFC3339),
			ScanType:  "SubdomainX Reconnaissance",
			TargetURL: "Multiple targets",
		},
		Sites: zapSites,
	}

	// Write XML file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create ZAP file: %v", err)
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("failed to encode ZAP XML: %v", err)
	}

	return nil
}

// extractHost extracts host from URL
func extractHost(url string) string {
	// Remove protocol
	if len(url) > 8 && url[:8] == "https://" {
		url = url[8:]
	} else if len(url) > 7 && url[:7] == "http://" {
		url = url[7:]
	}

	// Extract host (everything before first slash)
	for i, char := range url {
		if char == '/' {
			return url[:i]
		}
	}
	return url
}

// extractPortAndSSL extracts port and SSL info from URL
func extractPortAndSSL(url string) (string, string) {
	// Simple extraction - in production, use url.Parse
	if len(url) > 8 && url[:8] == "https://" {
		return "443", "true"
	}
	return "80", "false"
}
