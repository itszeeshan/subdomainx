package output

import (
	"fmt"
	"path/filepath"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/types"
)

// Generate creates output files based on the configuration and results
func Generate(cfg *config.Config, subdomainResults []types.SubdomainResult, httpResults []types.HTTPResult, portResults []types.PortResult) error {
	// Create scan results structure
	results := &types.ScanResults{
		Subdomains: subdomainResults,
		HTTP:       httpResults,
		Ports:      portResults,
	}

	// Generate output based on format
	switch cfg.OutputFormat {
	case "json":
		return generateJSON(cfg, results)
	case "txt":
		return generateTXT(cfg, results)
	case "html":
		return generateHTML(cfg, results)
	case "zap":
		return generateZAP(cfg, results)
	case "burp":
		return generateBurp(cfg, results)
	case "nessus":
		return generateNessus(cfg, results)
	case "csv":
		return generateCSV(cfg, results)
	default:
		return fmt.Errorf("unsupported output format: %s. Supported formats: json, txt, html, zap, burp, nessus, csv", cfg.OutputFormat)
	}
}

// generateJSON creates JSON output files
func generateJSON(cfg *config.Config, results *types.ScanResults) error {
	// Main results file
	mainFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_results.json", cfg.UniqueName))
	if err := WriteJSON(mainFile, results); err != nil {
		return fmt.Errorf("failed to write main JSON file: %v", err)
	}

	// Subdomains only file
	subdomainsFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_subdomains.json", cfg.UniqueName))
	if err := WriteJSON(subdomainsFile, results.Subdomains); err != nil {
		return fmt.Errorf("failed to write subdomains JSON file: %v", err)
	}

	// HTTP results file
	if len(results.HTTP) > 0 {
		httpFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_http.json", cfg.UniqueName))
		if err := WriteJSON(httpFile, results.HTTP); err != nil {
			return fmt.Errorf("failed to write HTTP JSON file: %v", err)
		}
	}

	// Port scan results file
	if len(results.Ports) > 0 {
		portsFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_ports.json", cfg.UniqueName))
		if err := WriteJSON(portsFile, results.Ports); err != nil {
			return fmt.Errorf("failed to write ports JSON file: %v", err)
		}
	}

	return nil
}

// generateTXT creates text output files
func generateTXT(cfg *config.Config, results *types.ScanResults) error {
	// Subdomains only file (most common use case)
	subdomainsFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_subdomains.txt", cfg.UniqueName))
	if err := WriteTXT(subdomainsFile, results.Subdomains); err != nil {
		return fmt.Errorf("failed to write subdomains TXT file: %v", err)
	}

	// HTTP results file
	if len(results.HTTP) > 0 {
		httpFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_http.txt", cfg.UniqueName))
		if err := WriteHTTPTXT(httpFile, results.HTTP); err != nil {
			return fmt.Errorf("failed to write HTTP TXT file: %v", err)
		}
	}

	// Port scan results file
	if len(results.Ports) > 0 {
		portsFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_ports.txt", cfg.UniqueName))
		if err := WritePortsTXT(portsFile, results.Ports); err != nil {
			return fmt.Errorf("failed to write ports TXT file: %v", err)
		}
	}

	return nil
}

// generateHTML creates HTML output files
func generateHTML(cfg *config.Config, results *types.ScanResults) error {
	// Main HTML report
	htmlFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_report.html", cfg.UniqueName))
	if err := WriteHTML(htmlFile, results); err != nil {
		return fmt.Errorf("failed to write HTML report: %v", err)
	}

	return nil
}

// generateZAP creates ZAP-compatible XML output files
func generateZAP(cfg *config.Config, results *types.ScanResults) error {
	// Main ZAP file
	zapFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_zap.xml", cfg.UniqueName))
	if err := WriteZAP(zapFile, results); err != nil {
		return fmt.Errorf("failed to write ZAP file: %v", err)
	}

	return nil
}

// generateBurp creates Burp Suite-compatible XML output files
func generateBurp(cfg *config.Config, results *types.ScanResults) error {
	// Main Burp file
	burpFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_burp.xml", cfg.UniqueName))
	if err := WriteBurp(burpFile, results); err != nil {
		return fmt.Errorf("failed to write Burp file: %v", err)
	}

	return nil
}

// generateNessus creates Nessus-compatible XML output files
func generateNessus(cfg *config.Config, results *types.ScanResults) error {
	// Main Nessus file
	nessusFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_nessus.xml", cfg.UniqueName))
	if err := WriteNessus(nessusFile, results); err != nil {
		return fmt.Errorf("failed to write Nessus file: %v", err)
	}

	return nil
}

// generateCSV creates CSV output files
func generateCSV(cfg *config.Config, results *types.ScanResults) error {
	// Main CSV file
	csvFile := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s_results.csv", cfg.UniqueName))
	if err := WriteCSV(csvFile, results); err != nil {
		return fmt.Errorf("failed to write CSV file: %v", err)
	}

	return nil
}
