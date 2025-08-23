package output

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/internal/types"
)

// WriteHTML creates an HTML report of the scan results
func WriteHTML(filename string, results *types.ScanResults) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write HTML header
	html := generateHTMLHeader()

	// Write summary section
	html += generateSummarySection(results)

	// Write subdomains section
	html += generateSubdomainsSection(results.Subdomains)

	// Write HTTP results section
	if len(results.HTTP) > 0 {
		html += generateHTTPSection(results.HTTP)
	}

	// Write port scan results section
	if len(results.Ports) > 0 {
		html += generatePortsSection(results.Ports)
	}

	// Write HTML footer
	html += generateHTMLFooter()

	_, err = file.WriteString(html)
	return err
}

func generateHTMLHeader() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SubdomainX Scan Report</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 2.5em;
            font-weight: 300;
        }
        .header p {
            margin: 10px 0 0 0;
            opacity: 0.9;
        }
        .content {
            padding: 30px;
        }
        .section {
            margin-bottom: 40px;
        }
        .section h2 {
            color: #333;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
            margin-bottom: 20px;
        }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .stat-card {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            text-align: center;
            border-left: 4px solid #667eea;
        }
        .stat-number {
            font-size: 2em;
            font-weight: bold;
            color: #667eea;
        }
        .stat-label {
            color: #666;
            margin-top: 5px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
            background: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        th, td {
            padding: 12px 15px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #667eea;
            color: white;
            font-weight: 600;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .status-200 { color: #28a745; }
        .status-301, .status-302 { color: #ffc107; }
        .status-404 { color: #dc3545; }
        .status-500 { color: #dc3545; }
        .source-tag {
            background: #e9ecef;
            padding: 2px 8px;
            border-radius: 12px;
            font-size: 0.8em;
            color: #495057;
        }
        .footer {
            background: #f8f9fa;
            padding: 20px;
            text-align: center;
            color: #666;
            border-top: 1px solid #dee2e6;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 SubdomainX Scan Report</h1>
            <p>Generated on ` + time.Now().Format("January 2, 2006 at 15:04:05") + `</p>
        </div>
        <div class="content">`
}

func generateSummarySection(results *types.ScanResults) string {
	return fmt.Sprintf(`
        <div class="section">
            <h2>📊 Scan Summary</h2>
            <div class="stats">
                <div class="stat-card">
                    <div class="stat-number">%d</div>
                    <div class="stat-label">Subdomains Found</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number">%d</div>
                    <div class="stat-label">HTTP Services</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number">%d</div>
                    <div class="stat-label">Hosts with Open Ports</div>
                </div>
            </div>
        </div>`, len(results.Subdomains), len(results.HTTP), len(results.Ports))
}

func generateSubdomainsSection(subdomains []types.SubdomainResult) string {
	if len(subdomains) == 0 {
		return `
        <div class="section">
            <h2>🌐 Subdomains</h2>
            <p>No subdomains found.</p>
        </div>`
	}

	html := `
        <div class="section">
            <h2>🌐 Subdomains Found</h2>
            <table>
                <thead>
                    <tr>
                        <th>Subdomain</th>
                        <th>Source</th>
                        <th>IP Addresses</th>
                    </tr>
                </thead>
                <tbody>`

	for _, subdomain := range subdomains {
		ips := strings.Join(subdomain.IPs, ", ")
		if ips == "" {
			ips = "N/A"
		}
		html += fmt.Sprintf(`
                    <tr>
                        <td><strong>%s</strong></td>
                        <td><span class="source-tag">%s</span></td>
                        <td>%s</td>
                    </tr>`, subdomain.Subdomain, subdomain.Source, ips)
	}

	html += `
                </tbody>
            </table>
        </div>`

	return html
}

func generateHTTPSection(httpResults []types.HTTPResult) string {
	html := `
        <div class="section">
            <h2>🌍 HTTP Services</h2>
            <table>
                <thead>
                    <tr>
                        <th>URL</th>
                        <th>Status</th>
                        <th>Title</th>
                        <th>Content Length</th>
                        <th>Technologies</th>
                    </tr>
                </thead>
                <tbody>`

	for _, result := range httpResults {
		statusClass := fmt.Sprintf("status-%d", result.StatusCode)
		technologies := strings.Join(result.Technologies, ", ")
		if technologies == "" {
			technologies = "N/A"
		}

		html += fmt.Sprintf(`
                    <tr>
                        <td><a href="%s" target="_blank">%s</a></td>
                        <td class="%s">%d</td>
                        <td>%s</td>
                        <td>%d</td>
                        <td>%s</td>
                    </tr>`, result.URL, result.URL, statusClass, result.StatusCode, result.Title, result.ContentLength, technologies)
	}

	html += `
                </tbody>
            </table>
        </div>`

	return html
}

func generatePortsSection(portResults []types.PortResult) string {
	html := `
        <div class="section">
            <h2>🔌 Port Scan Results</h2>
            <table>
                <thead>
                    <tr>
                        <th>Host</th>
                        <th>IP</th>
                        <th>Port</th>
                        <th>Protocol</th>
                        <th>State</th>
                        <th>Service</th>
                        <th>Version</th>
                    </tr>
                </thead>
                <tbody>`

	for _, result := range portResults {
		for _, port := range result.Ports {
			html += fmt.Sprintf(`
                    <tr>
                        <td><strong>%s</strong></td>
                        <td>%s</td>
                        <td>%d</td>
                        <td>%s</td>
                        <td>%s</td>
                        <td>%s</td>
                        <td>%s</td>
                    </tr>`, result.Host, result.IP, port.Number, port.Protocol, port.State, port.Service, port.Version)
		}
	}

	html += `
                </tbody>
            </table>
        </div>`

	return html
}

func generateHTMLFooter() string {
	return `
        </div>
        <div class="footer">
            <p>Report generated by SubdomainX - All-in-one subdomain enumeration tool</p>
        </div>
    </div>
</body>
</html>`
}
