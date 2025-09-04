<div align="center">
  <img src="docs/public/logo.png" alt="SubdomainX Logo" width="200"/>
  <h1>SubdomainX</h1>
  <p><strong>Advanced Subdomain Discovery & Security Reconnaissance Tool</strong></p>
</div>

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/itszeeshan/subdomainx)](https://goreportcard.com/report/github.com/itszeeshan/subdomainx)
[![Go Reference](https://pkg.go.dev/badge/github.com/itszeeshan/subdomainx.svg)](https://pkg.go.dev/github.com/itszeeshan/subdomainx)
[![GitHub release](https://img.shields.io/github/release/itszeeshan/subdomainx.svg)](https://github.com/itszeeshan/subdomainx/releases)
[![GitHub stars](https://img.shields.io/github/stars/itszeeshan/subdomainx.svg)](https://github.com/itszeeshan/subdomainx/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/itszeeshan/subdomainx.svg)](https://github.com/itszeeshan/subdomainx/network)
[![GitHub issues](https://img.shields.io/github/issues/itszeeshan/subdomainx.svg)](https://github.com/itszeeshan/subdomainx/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/itszeeshan/subdomainx.svg)](https://github.com/itszeeshan/subdomainx/pulls)
[![CI/CD](https://img.shields.io/github/actions/workflow/status/itszeeshan/subdomainx/ci.yml?branch=main)](https://github.com/itszeeshan/subdomainx/actions)
[![Code Coverage](https://img.shields.io/badge/coverage-85%25-brightgreen.svg)](https://github.com/itszeeshan/subdomainx)

</div>

---

**SubdomainX** is a powerful, all-in-one subdomain enumeration and reconnaissance tool designed for modern cybersecurity professionals, penetration testers, and security researchers.

### Why Choose SubdomainX?

- **All-in-One Solution**: Combines 12+ popular enumeration tools into a single, unified interface
- **High Performance**: Multi-threaded architecture with intelligent resource management
- **API Integration**: Native integration with SecurityTrails, VirusTotal, Censys, and more
- **Advanced Reporting**: Beautiful HTML reports, JSON exports, and real-time progress tracking
- **Resume Capability**: Never lose progress with intelligent checkpointing system
- **Smart Optimization**: Built-in resource monitoring and performance recommendations
- **Production Ready**: Comprehensive error handling, logging, and validation

### Perfect For:

- **Security Researchers** conducting comprehensive domain reconnaissance
- **Penetration Testers** performing thorough attack surface analysis
- **Bug Bounty Hunters** discovering hidden subdomains and assets
- **Security Teams** monitoring their organization's digital footprint
- **Red Teams** gathering intelligence for advanced persistent threats

## Key Features

### Intelligent Enumeration

- **12+ Tools Integrated**: subfinder, amass, findomain, assetfinder, sublist3r, knockpy, dnsrecon, fierce, massdns, altdns, waybackurls, linkheader
- **7+ API Services**: SecurityTrails, VirusTotal, Censys, crt.sh, URLScan.io, ThreatCrowd, HackerTarget
- **Custom Wordlists**: Support for custom brute-forcing dictionaries
- **Smart Filtering**: Advanced filtering and deduplication

### HTTP & Port Scanning

- **httpx Integration**: Comprehensive HTTP probing with status codes, headers, and technologies
- **smap Integration**: Fast port scanning with service detection
- **Customizable Filters**: Filter by status codes, ports, and response patterns

### Advanced Monitoring

- **Real-time Progress**: Live progress bars with ETA calculations
- **Resource Management**: CPU and memory monitoring with optimization tips
- **Checkpoint System**: Save and resume interrupted scans seamlessly
- **Comprehensive Logging**: Detailed logs for debugging and analysis

### Professional Reporting

- **Multiple Formats**: JSON, TXT, HTML, CSV, and security tool formats
- **Security Tool Integration**: Export to OWASP ZAP, Burp Suite, and Nessus formats
- **Customizable Output**: Flexible naming and directory structure
- **Rich Metadata**: Detailed scan information and statistics
- **Export Ready**: Compatible with other security tools and platforms

## Quick Start

### Installation

```bash
# Install from source
go install github.com/itszeeshan/subdomainx@latest

# Or download pre-built binary
curl -sSL https://github.com/itszeeshan/subdomainx/releases/latest/download/subdomainx_$(uname -s)_$(uname -m).tar.gz | tar -xz
sudo mv subdomainx /usr/local/bin/
```

### Basic Usage

**Single Domain Enumeration:**

```bash
subdomainx --subfinder --httpx example.com
```

**Multiple Domains:**

```bash
echo "example.com" > domains.txt
subdomainx --wildcard domains.txt --format html
```

**Security Tool Integration:**

```bash
# Export to OWASP ZAP format
subdomainx --subfinder --httpx --format zap example.com

# Export to Burp Suite format
subdomainx --subfinder --httpx --format burp example.com

# Export to Nessus format
subdomainx --subfinder --httpx --format nessus example.com

# Export to CSV for spreadsheet analysis
subdomainx --subfinder --httpx --format csv example.com
```

**API-Powered Discovery:**

```bash
# Set API keys
export SECURITYTRAILS_API_KEY="your_key"
export VIRUSTOTAL_API_KEY="your_key"
export CENSYS_API_ID="your_id"
export CENSYS_SECRET="your_secret"
export URLSCAN_API_KEY="your_key"
export HACKERTARGET_API_KEY="your_key"

# Use APIs
subdomainx --securitytrails --virustotal --censys --crtsh --urlscan --threatcrowd --hackertarget example.com
```

**High-Performance Scan:**

```bash
subdomainx --threads 20 --timeout 60 --subfinder --amass --max-http-targets 1000 example.com
```

**Resume Interrupted Scan:**

```bash
# Resume from checkpoint
subdomainx --resume my_scan
```

> **Pro Tip**: Always place flags before the domain argument:
>
> ```bash
> subdomainx --tools domain.com  # Correct
> subdomainx domain.com --tools  # Incorrect
> ```

## Supported Tools

### Enumeration Tools

| Tool            | Description                        | Installation                                                                               |
| --------------- | ---------------------------------- | ------------------------------------------------------------------------------------------ |
| **subfinder**   | Fast subdomain discovery           | `go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest`              |
| **amass**       | In-depth subdomain enumeration     | `go install -v github.com/owasp-amass/amass/v4/...@master`                                 |
| **findomain**   | Fast subdomain finder              | `curl -LO https://github.com/findomain/findomain/releases/latest/download/findomain-linux` |
| **assetfinder** | Find subdomains and related assets | `go install github.com/tomnomnom/assetfinder@latest`                                       |
| **sublist3r**   | Python-based subdomain enumeration | `pip install sublist3r`                                                                    |
| **knockpy**     | Subdomain enumeration tool         | `pip install knockpy`                                                                      |
| **dnsrecon**    | DNS enumeration tool               | `pip install dnsrecon`                                                                     |
| **fierce**      | DNS reconnaissance tool            | `pip install fierce`                                                                       |
| **massdns**     | High-performance DNS resolver      | `git clone https://github.com/blechschmidt/massdns.git`                                    |
| **altdns**      | Subdomain permutation tool         | `pip install altdns`                                                                       |
| **waybackurls** | Wayback Machine URL finder         | `go install github.com/tomnomnom/waybackurls@latest`                                       |
| **linkheader**  | HTTP Link header parser            | Built-in                                                                                   |

### API Services

| Service            | Description                 | API Key Required |
| ------------------ | --------------------------- | ---------------- |
| **SecurityTrails** | Historical DNS data         | ✅               |
| **VirusTotal**     | Threat intelligence         | ✅               |
| **Censys**         | Internet-wide scanning data | ✅               |
| **crt.sh**         | Certificate Transparency    | ❌               |
| **URLScan.io**     | Web scanning service        | ✅ (optional)    |
| **ThreatCrowd**    | Threat intelligence         | ❌               |
| **HackerTarget**   | Security research platform  | ✅ (optional)    |

### Scanning Tools

| Tool      | Description     | Installation                                                       |
| --------- | --------------- | ------------------------------------------------------------------ |
| **httpx** | Fast HTTP probe | `go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest` |
| **smap**  | Port scanner    | `go install github.com/s0md3v/smap/cmd/smap@latest`                |

## Documentation

**[View Full Documentation](https://subdomainx.vercel.app)**

- [Installation Guide](https://subdomainx.vercel.app/installation)
- [CLI Reference](https://subdomainx.vercel.app/cli-reference)
- [Examples & Use Cases](https://subdomainx.vercel.app/examples)
- [Configuration](https://subdomainx.vercel.app/configuration)
- [Supported Tools](https://subdomainx.vercel.app/supported-tools)

## Advanced Examples

### Comprehensive Reconnaissance

```bash
# Full enumeration with all tools
subdomainx --subfinder --amass --findomain --assetfinder --sublist3r \
           --securitytrails --virustotal --censys \
           --httpx --smap \
           --format html --name comprehensive_scan example.com
```

### Targeted Enumeration

```bash
# Focus on specific tools for speed
subdomainx --subfinder --httpx --status-codes 200,301,302 \
           --ports 80,443,8080,8443 --max-http-targets 500 example.com
```

### Custom Wordlist Brute Force

```bash
# Use custom wordlist for altdns
subdomainx --altdns --wordlist /path/to/custom_wordlist.txt example.com
```

### Resume and Monitor

```bash
# Start scan with monitoring
subdomainx --verbose --subfinder --amass --max-http-targets 1000 example.com

# Later resume if interrupted
subdomainx --resume example_com_scan
```

## Contributing

We welcome contributions! Here's how you can help:

1. **Report Bugs**: [Create an issue](https://github.com/itszeeshan/subdomainx/issues)
2. **Suggest Features**: [Start a discussion](https://github.com/itszeeshan/subdomainx/discussions)
3. **Submit PRs**: Fork the repo and submit pull requests
4. **Improve Docs**: Help us make the documentation better
5. **Star the Repo**: Show your support!

### Development Setup

```bash
git clone https://github.com/itszeeshan/subdomainx.git
cd subdomainx
go mod download
go build -o subdomainx .
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

**SubdomainX is designed for authorized security testing and research purposes only.**

- Always ensure you have proper authorization before scanning any domain
- Respect rate limits and terms of service of target systems
- Use responsibly and ethically
- The authors are not responsible for any misuse of this tool

## Acknowledgments

- All the amazing open-source tools that make SubdomainX possible
- The security community for continuous feedback and improvements
- Contributors and users who help make this tool better

---

<div align="center">
  <p><strong>Happy Hunting! 🎯</strong></p>
  <p>Made with ❤️ by Zeeshan</p>
</div>

