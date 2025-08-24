package output

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/itszeeshan/subdomainx/internal/types"
)

// BurpItem represents an item in Burp Suite XML format
type BurpItem struct {
	XMLName     xml.Name `xml:"item"`
	Time        string   `xml:"time"`
	URL         string   `xml:"url"`
	Host        string   `xml:"host"`
	Port        string   `xml:"port"`
	Protocol    string   `xml:"protocol"`
	Method      string   `xml:"method"`
	Path        string   `xml:"path"`
	Extension   string   `xml:"extension"`
	Request     string   `xml:"request"`
	Status      string   `xml:"status"`
	Response    string   `xml:"response"`
	ResponseURL string   `xml:"responseRedirectUrl"`
	Comments    string   `xml:"comments"`
}

// BurpReport represents the complete Burp Suite XML report
type BurpReport struct {
	XMLName xml.Name   `xml:"items"`
	Items   []BurpItem `xml:"item"`
}

// WriteBurp creates a Burp Suite-compatible XML file
func WriteBurp(filename string, results *types.ScanResults) error {
	var burpItems []BurpItem

	// Convert HTTP results to Burp items
	for _, http := range results.HTTP {
		host, port, protocol, path := parseURL(http.URL)

		item := BurpItem{
			Time:        time.Now().Format(time.RFC3339),
			URL:         http.URL,
			Host:        host,
			Port:        port,
			Protocol:    protocol,
			Method:      "GET",
			Path:        path,
			Extension:   getExtension(path),
			Request:     generateRequest(http.URL),
			Status:      fmt.Sprintf("%d", http.StatusCode),
			Response:    generateResponse(http),
			ResponseURL: http.URL,
			Comments:    generateComments(http),
		}
		burpItems = append(burpItems, item)
	}

	// Create Burp report
	report := BurpReport{
		Items: burpItems,
	}

	// Write XML file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create Burp file: %v", err)
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("failed to encode Burp XML: %v", err)
	}

	return nil
}

// parseURL parses URL into components
func parseURL(url string) (host, port, protocol, path string) {
	// Simple URL parsing - in production, use url.Parse
	if len(url) > 8 && url[:8] == "https://" {
		protocol = "https"
		port = "443"
		url = url[8:]
	} else if len(url) > 7 && url[:7] == "http://" {
		protocol = "http"
		port = "80"
		url = url[7:]
	}

	// Extract host and path
	for i, char := range url {
		if char == '/' {
			host = url[:i]
			path = url[i:]
			return
		}
	}
	host = url
	path = "/"
	return
}

// getExtension extracts file extension from path
func getExtension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			break
		}
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return ""
}

// generateRequest creates a simple HTTP request
func generateRequest(url string) string {
	host, _, _, path := parseURL(url)
	return fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: SubdomainX/1.0\r\n\r\n", path, host)
}

// generateResponse creates a simple HTTP response
func generateResponse(http types.HTTPResult) string {
	statusText := "OK"
	if http.StatusCode >= 400 {
		statusText = "Not Found"
	}

	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", http.StatusCode, statusText)
	response += "Content-Type: text/html\r\n"
	if http.ContentLength > 0 {
		response += fmt.Sprintf("Content-Length: %d\r\n", http.ContentLength)
	}
	response += "\r\n"

	if http.Title != "" {
		response += fmt.Sprintf("<html><head><title>%s</title></head><body></body></html>", http.Title)
	}

	return response
}

// generateComments creates comments from HTTP result metadata
func generateComments(http types.HTTPResult) string {
	var comments []string

	if http.Title != "" {
		comments = append(comments, "Title: "+http.Title)
	}

	if len(http.Technologies) > 0 {
		comments = append(comments, "Technologies: "+fmt.Sprintf("%v", http.Technologies))
	}

	if http.ContentLength > 0 {
		comments = append(comments, fmt.Sprintf("Content-Length: %d", http.ContentLength))
	}

	if len(comments) > 0 {
		return fmt.Sprintf("SubdomainX: %s", fmt.Sprintf("%s", comments))
	}

	return "SubdomainX: Discovered subdomain"
}
