package output

import (
	"encoding/json"
	"os"

	"github.com/itszeeshan/subdomainx/internal/types"
)

// WriteJSON writes data to a JSON file
func WriteJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// WriteSubdomainJSON writes subdomain results to JSON
func WriteSubdomainJSON(filename string, subdomains []types.SubdomainResult) error {
	return WriteJSON(filename, subdomains)
}

// WriteHTTPJSON writes HTTP results to JSON
func WriteHTTPJSON(filename string, httpResults []types.HTTPResult) error {
	return WriteJSON(filename, httpResults)
}

// WritePortsJSON writes port scan results to JSON
func WritePortsJSON(filename string, portResults []types.PortResult) error {
	return WriteJSON(filename, portResults)
}

// WriteScanResultsJSON writes complete scan results to JSON
func WriteScanResultsJSON(filename string, results *types.ScanResults) error {
	return WriteJSON(filename, results)
}
