# SubdomainX 🔍

SubdomainX is an all-in-one subdomain enumeration tool that combines multiple popular subdomain discovery tools into a single, easy-to-use interface. It provides comprehensive subdomain enumeration, HTTP scanning, and port scanning capabilities with beautiful HTML reports.

## Features ✨

- **Multiple Enumeration Tools**: Integrates with popular tools like subfinder, amass, findomain, assetfinder, and more
- **HTTP Scanning**: Discovers web services and extracts titles, status codes, and technologies
- **Port Scanning**: Identifies open ports and services on discovered hosts
- **Concurrent Processing**: Multi-threaded execution for faster results
- **Multiple Output Formats**: JSON, TXT, and beautiful HTML reports
- **DNS Caching**: Efficient DNS resolution with caching
- **Rate Limiting**: Configurable rate limiting to avoid overwhelming targets
- **Retry Mechanism**: Automatic retry with exponential backoff
- **Filtering**: Filter results by status codes, ports, and other criteria

## Supported Tools 🛠️

### Subdomain Enumeration

- **subfinder** - Fast subdomain discovery
- **amass** - Comprehensive subdomain enumeration
- **findomain** - Fast and cross-platform subdomain discovery
- **assetfinder** - Find subdomains related to a domain
- **sublist3r** - Subdomain enumeration using OSINT
- **knockpy** - Subdomain enumeration tool
- **dnsrecon** - DNS enumeration and reconnaissance
- **fierce** - DNS reconnaissance tool
- **massdns** - High-performance DNS stub resolver
- **altdns** - Subdomain permutation and alteration

### Scanning Tools

- **httpx** - Fast and multi-purpose HTTP probe
- **smap** - Port scanner and service discovery

## Installation 📦

### Prerequisites

Make sure you have Go 1.21+ installed and the following tools available in your PATH:

```bash
# Install required tools (examples for different systems)
# Ubuntu/Debian
sudo apt install subfinder amass findomain assetfinder sublist3r knockpy dnsrecon fierce massdns altdns httpx smap

# macOS
brew install subfinder amass findomain assetfinder sublist3r knockpy dnsrecon fierce massdns altdns httpx smap

# Or install individually from their respective repositories
```

### Quick Install

Install SubdomainX directly using Go:

```bash
# Install the latest version
go install github.com/itszeeshan/subdomainx/cmd/subdomainx@latest

# (Optional) Make it globally accessible
# The binary will be installed to $(go env GOPATH)/bin/subdomainx
# Make sure $(go env GOPATH)/bin is in your PATH, or move it to /usr/local/bin:
sudo mv $(go env GOPATH)/bin/subdomainx /usr/local/bin/

# Now you can run subdomainx from anywhere 🚀
subdomainx --help
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/itszeeshan/subdomainx.git
cd subdomainx

# Build the binary
go build -o subdomainx ./cmd/subdomainx

# Or use the Makefile
make build

# Or install globally
go install ./cmd/subdomainx
```

### Tool Installation Check

SubdomainX includes built-in tool checking and installation assistance:

```bash
# Check which tools are available
./subdomainx --check-tools

# Get installation instructions for missing tools
./subdomainx --install-tools

# The tool will automatically detect missing tools and offer to help install them
```

## Usage 🚀

### Basic Usage

1. **Create a domains file** with your target domains:

   ```bash
   echo "example.com" > domains.txt
   echo "test.org" >> domains.txt
   ```

2. **Run the scan**:

   ```bash
   # Basic scan with all available tools
   ./subdomainx --wildcard domains.txt

   # Use only specific tools
   ./subdomainx --wildcard domains.txt --amass --subfinder --httpx

   # Generate HTML report
   ./subdomainx --wildcard domains.txt --format html --name my_scan
   ```

### Configuration (Optional)

The tool can use a YAML configuration file located at `configs/default.yaml` for default settings. All settings can also be specified via CLI flags, which take precedence over the config file:

```yaml
wildcard_file: domains.txt # File containing target domains
unique_name: scan # Unique name for output files
output_dir: output # Output directory
output_format: json # Output format: json, txt, html
threads: 10 # Number of concurrent threads
retries: 3 # Number of retry attempts
timeout: 30 # Timeout in seconds
rate_limit: 100 # Rate limit per second
wordlist: wordlists/default.txt # Custom wordlist (optional)
filters:
  status_code: "200,301,302" # Filter by HTTP status codes
  ports: "80,443,8080" # Filter by ports
tools:
  subfinder: true # Enable/disable specific tools
  findomain: true
  assetfinder: true
  amass: true
  sublist3r: true
  knockpy: true
  dnsrecon: true
  fierce: true
  massdns: true
  altdns: true
```

### Command Line Options

```bash
# Required
--wildcard FILE        Path to file containing target domains (one per line)

# Tool Selection (use specific tools, otherwise all available)
--subfinder            Use subfinder tool
--amass                Use amass tool
--findomain            Use findomain tool
--assetfinder          Use assetfinder tool
--sublist3r            Use sublist3r tool
--knockpy              Use knockpy tool
--dnsrecon             Use dnsrecon tool
--fierce               Use fierce tool
--massdns              Use massdns tool
--altdns               Use altdns tool
--httpx                Use httpx for HTTP scanning
--smap                 Use smap for port scanning

# Output Options
--name NAME            Unique name for output files (default: scan)
--format FORMAT        Output format: json, txt, html (default: json)
--output DIR           Output directory (default: output)

# Performance Options
--threads N            Number of concurrent threads (default: 10)
--retries N            Number of retry attempts (default: 3)
--timeout N            Timeout in seconds (default: 30)
--rate-limit N         Rate limit per second (default: 100)

# Utility Options
--help                 Show help message
--version              Show version information
--check-tools          Check tool availability
--install-tools        Show installation instructions
--config FILE          Use custom configuration file (optional)
-v, --verbose          Enable verbose output
```

## Output 📊

The tool generates multiple output files based on your configuration:

### JSON Output

- `{name}_results.json` - Complete scan results
- `{name}_subdomains.json` - Subdomain enumeration results
- `{name}_http.json` - HTTP scanning results
- `{name}_ports.json` - Port scanning results

### TXT Output

- `{name}_subdomains.txt` - List of discovered subdomains
- `{name}_http.txt` - HTTP results in tab-separated format
- `{name}_ports.txt` - Port scan results in tab-separated format

### HTML Output

- `{name}_report.html` - Beautiful HTML report with:
  - Scan summary with statistics
  - Subdomain discovery results
  - HTTP service details
  - Port scan results
  - Interactive tables and filtering

## Examples 💡

### Basic Subdomain Enumeration

```bash
# Create domains file
echo "example.com" > domains.txt

# Run scan with all available tools
./subdomainx --wildcard domains.txt

# Check results
ls output/
cat output/scan_subdomains.txt
```

### Custom Configuration

```yaml
# configs/custom.yaml
wildcard_file: targets.txt
unique_name: penetration_test
output_format: html
threads: 20
tools:
  subfinder: true
  amass: true
  findomain: true
  httpx: true
```

### Advanced Usage

```bash
# Run with specific tools only
./subdomainx --wildcard domains.txt --amass --subfinder --httpx

# Generate HTML report
./subdomainx --wildcard domains.txt --format html --name my_scan

# High-performance scan
./subdomainx --wildcard domains.txt --threads 20 --timeout 60

# Use custom configuration file
./subdomainx --wildcard domains.txt --config custom.yaml
```

## Contributing 🤝

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/subdomainx.git
cd subdomainx

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build
go build ./cmd/subdomainx
```

## License 📄

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer ⚠️

This tool is designed for authorized security testing and research purposes only. Always ensure you have proper authorization before scanning any domain or network. The authors are not responsible for any misuse of this tool.

## Acknowledgments 🙏

This tool integrates with and builds upon the excellent work of the following open-source projects:

- [subfinder](https://github.com/projectdiscovery/subfinder)
- [amass](https://github.com/owasp-amass/amass)
- [findomain](https://github.com/findomain/findomain)
- [assetfinder](https://github.com/tomnomnom/assetfinder)
- [httpx](https://github.com/projectdiscovery/httpx)
- [smap](https://github.com/s0md3v/Smap)
- And many others...

## Support 💬

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/itszeeshan/subdomainx/issues) page
2. Create a new issue with detailed information

---

**Happy Hunting! 🎯**
