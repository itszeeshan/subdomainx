package output

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/internal/types"
)

const itemsPerPage = 50

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
            background: linear-gradient(135deg, #16a34a 0%, #15803d 100%);
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
            border-bottom: 2px solid #16a34a;
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
            border-left: 4px solid #16a34a;
        }
        .stat-number {
            font-size: 2em;
            font-weight: bold;
            color: #16a34a;
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
            background-color: #16a34a;
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
        .pagination {
            display: flex;
            justify-content: center;
            align-items: center;
            margin-top: 20px;
            gap: 10px;
        }
        .pagination button {
            padding: 8px 16px;
            border: 1px solid #ddd;
            background: white;
            cursor: pointer;
            border-radius: 4px;
            transition: all 0.3s ease;
        }
        .pagination button:hover {
            background: #16a34a;
            color: white;
            border-color: #16a34a;
        }
        .pagination button:disabled {
            background: #f5f5f5;
            color: #999;
            cursor: not-allowed;
        }
        .pagination button.active {
            background: #16a34a;
            color: white;
            border-color: #16a34a;
        }
        .pagination-info {
            margin: 0 20px;
            color: #666;
        }
        .table-container {
            position: relative;
        }
        .loading {
            display: none;
            text-align: center;
            padding: 20px;
            color: #666;
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

	// Generate all rows as JSON for JavaScript pagination
	var rows []string
	for _, subdomain := range subdomains {
		ips := strings.Join(subdomain.IPs, ", ")
		if ips == "" {
			ips = "N/A"
		}
		row := fmt.Sprintf(`{"subdomain": "%s", "source": "%s", "ips": "%s"}`,
			subdomain.Subdomain, subdomain.Source, ips)
		rows = append(rows, row)
	}
	rowsJSON := "[" + strings.Join(rows, ",") + "]"

	html := fmt.Sprintf(`
        <div class="section">
            <h2>🌐 Subdomains Found (%d total)</h2>
            <div class="table-container">
                <table id="subdomains-table">
                    <thead>
                        <tr>
                            <th>Subdomain</th>
                            <th>Source</th>
                            <th>IP Addresses</th>
                        </tr>
                    </thead>
                    <tbody id="subdomains-tbody">
                    </tbody>
                </table>
                <div class="loading" id="subdomains-loading">Loading...</div>
                <div class="pagination" id="subdomains-pagination">
                    <button onclick="changePage('subdomains', -1)">Previous</button>
                    <span class="pagination-info" id="subdomains-info">Page 1 of 1</span>
                    <button onclick="changePage('subdomains', 1)">Next</button>
                </div>
            </div>
        </div>
        <script>
            const subdomainsData = %s;
            let subdomainsCurrentPage = 1;
            const subdomainsPerPage = %d;
            
            function renderSubdomainsTable() {
                const tbody = document.getElementById('subdomains-tbody');
                const startIndex = (subdomainsCurrentPage - 1) * subdomainsPerPage;
                const endIndex = startIndex + subdomainsPerPage;
                const pageData = subdomainsData.slice(startIndex, endIndex);
                
                tbody.innerHTML = '';
                pageData.forEach(item => {
                    const row = document.createElement('tr');
                    row.innerHTML = '<td><strong>' + item.subdomain + '</strong></td>' +
                                   '<td><span class="source-tag">' + item.source + '</span></td>' +
                                   '<td>' + item.ips + '</td>';
                    tbody.appendChild(row);
                });
                
                // Update pagination info
                const totalPages = Math.ceil(subdomainsData.length / subdomainsPerPage);
                document.getElementById('subdomains-info').textContent = 'Page ' + subdomainsCurrentPage + ' of ' + totalPages + ' (' + subdomainsData.length + ' total items)';
                
                // Update button states
                const prevBtn = document.querySelector('#subdomains-pagination button:first-child');
                const nextBtn = document.querySelector('#subdomains-pagination button:last-child');
                prevBtn.disabled = subdomainsCurrentPage === 1;
                nextBtn.disabled = subdomainsCurrentPage === totalPages;
            }
            
            function changePage(tableType, direction) {
                if (tableType === 'subdomains') {
                    const totalPages = Math.ceil(subdomainsData.length / subdomainsPerPage);
                    const newPage = subdomainsCurrentPage + direction;
                    if (newPage >= 1 && newPage <= totalPages) {
                        subdomainsCurrentPage = newPage;
                        renderSubdomainsTable();
                    }
                }
            }
            
            // Initialize table
            renderSubdomainsTable();
        </script>`, len(subdomains), rowsJSON, itemsPerPage)

	return html
}

func generateHTTPSection(httpResults []types.HTTPResult) string {
	if len(httpResults) == 0 {
		return `
        <div class="section">
            <h2>🌍 HTTP Services</h2>
            <p>No HTTP services found.</p>
        </div>`
	}

	// Generate all rows as JSON for JavaScript pagination
	var rows []string
	for _, result := range httpResults {
		statusClass := fmt.Sprintf("status-%d", result.StatusCode)
		technologies := strings.Join(result.Technologies, ", ")
		if technologies == "" {
			technologies = "N/A"
		}
		row := fmt.Sprintf(`{"url": "%s", "status": %d, "statusClass": "%s", "title": "%s", "contentLength": %d, "technologies": "%s"}`,
			result.URL, result.StatusCode, statusClass, result.Title, result.ContentLength, technologies)
		rows = append(rows, row)
	}
	rowsJSON := "[" + strings.Join(rows, ",") + "]"

	html := fmt.Sprintf(`
        <div class="section">
            <h2>🌍 HTTP Services (%d total)</h2>
            <div class="table-container">
                <table id="http-table">
                    <thead>
                        <tr>
                            <th>URL</th>
                            <th>Status</th>
                            <th>Title</th>
                            <th>Content Length</th>
                            <th>Technologies</th>
                        </tr>
                    </thead>
                    <tbody id="http-tbody">
                    </tbody>
                </table>
                <div class="loading" id="http-loading">Loading...</div>
                <div class="pagination" id="http-pagination">
                    <button onclick="changeHTTPPage(-1)">Previous</button>
                    <span class="pagination-info" id="http-info">Page 1 of 1</span>
                    <button onclick="changeHTTPPage(1)">Next</button>
                </div>
            </div>
        </div>
        <script>
            const httpData = %s;
            let httpCurrentPage = 1;
            const httpPerPage = %d;
            
            function renderHTTPTable() {
                const tbody = document.getElementById('http-tbody');
                const startIndex = (httpCurrentPage - 1) * httpPerPage;
                const endIndex = startIndex + httpPerPage;
                const pageData = httpData.slice(startIndex, endIndex);
                
                tbody.innerHTML = '';
                pageData.forEach(item => {
                    const row = document.createElement('tr');
                    row.innerHTML = '<td><a href="' + item.url + '" target="_blank">' + item.url + '</a></td>' +
                                   '<td class="' + item.statusClass + '">' + item.status + '</td>' +
                                   '<td>' + item.title + '</td>' +
                                   '<td>' + item.contentLength + '</td>' +
                                   '<td>' + item.technologies + '</td>';
                    tbody.appendChild(row);
                });
                
                // Update pagination info
                const totalPages = Math.ceil(httpData.length / httpPerPage);
                document.getElementById('http-info').textContent = 'Page ' + httpCurrentPage + ' of ' + totalPages + ' (' + httpData.length + ' total items)';
                
                // Update button states
                const prevBtn = document.querySelector('#http-pagination button:first-child');
                const nextBtn = document.querySelector('#http-pagination button:last-child');
                prevBtn.disabled = httpCurrentPage === 1;
                nextBtn.disabled = httpCurrentPage === totalPages;
            }
            
            function changeHTTPPage(direction) {
                const totalPages = Math.ceil(httpData.length / httpPerPage);
                const newPage = httpCurrentPage + direction;
                if (newPage >= 1 && newPage <= totalPages) {
                    httpCurrentPage = newPage;
                    renderHTTPTable();
                }
            }
            
            // Initialize table
            renderHTTPTable();
        </script>`, len(httpResults), rowsJSON, itemsPerPage)

	return html
}

func generatePortsSection(portResults []types.PortResult) string {
	if len(portResults) == 0 {
		return `
        <div class="section">
            <h2>🔌 Port Scan Results</h2>
            <p>No port scan results found.</p>
        </div>`
	}

	// Generate all rows as JSON for JavaScript pagination
	var rows []string
	for _, result := range portResults {
		for _, port := range result.Ports {
			row := fmt.Sprintf(`{"host": "%s", "ip": "%s", "port": %d, "protocol": "%s", "state": "%s", "service": "%s", "version": "%s"}`,
				result.Host, result.IP, port.Number, port.Protocol, port.State, port.Service, port.Version)
			rows = append(rows, row)
		}
	}
	rowsJSON := "[" + strings.Join(rows, ",") + "]"

	html := fmt.Sprintf(`
        <div class="section">
            <h2>🔌 Port Scan Results (%d total)</h2>
            <div class="table-container">
                <table id="ports-table">
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
                    <tbody id="ports-tbody">
                    </tbody>
                </table>
                <div class="loading" id="ports-loading">Loading...</div>
                <div class="pagination" id="ports-pagination">
                    <button onclick="changePortsPage(-1)">Previous</button>
                    <span class="pagination-info" id="ports-info">Page 1 of 1</span>
                    <button onclick="changePortsPage(1)">Next</button>
                </div>
            </div>
        </div>
        <script>
            const portsData = %s;
            let portsCurrentPage = 1;
            const portsPerPage = %d;
            
            function renderPortsTable() {
                const tbody = document.getElementById('ports-tbody');
                const startIndex = (portsCurrentPage - 1) * portsPerPage;
                const endIndex = startIndex + portsPerPage;
                const pageData = portsData.slice(startIndex, endIndex);
                
                tbody.innerHTML = '';
                pageData.forEach(item => {
                    const row = document.createElement('tr');
                    row.innerHTML = '<td><strong>' + item.host + '</strong></td>' +
                                   '<td>' + item.ip + '</td>' +
                                   '<td>' + item.port + '</td>' +
                                   '<td>' + item.protocol + '</td>' +
                                   '<td>' + item.state + '</td>' +
                                   '<td>' + item.service + '</td>' +
                                   '<td>' + item.version + '</td>';
                    tbody.appendChild(row);
                });
                
                // Update pagination info
                const totalPages = Math.ceil(portsData.length / portsPerPage);
                document.getElementById('ports-info').textContent = 'Page ' + portsCurrentPage + ' of ' + totalPages + ' (' + portsData.length + ' total items)';
                
                // Update button states
                const prevBtn = document.querySelector('#ports-pagination button:first-child');
                const nextBtn = document.querySelector('#ports-pagination button:last-child');
                prevBtn.disabled = portsCurrentPage === 1;
                nextBtn.disabled = portsCurrentPage === totalPages;
            }
            
            function changePortsPage(direction) {
                const totalPages = Math.ceil(portsData.length / portsPerPage);
                const newPage = portsCurrentPage + direction;
                if (newPage >= 1 && newPage <= totalPages) {
                    portsCurrentPage = newPage;
                    renderPortsTable();
                }
            }
            
            // Initialize table
            renderPortsTable();
        </script>`, len(rows), rowsJSON, itemsPerPage)

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
