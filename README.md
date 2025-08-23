# SubdomainX 🔍

All-in-one subdomain enumeration tool that combines multiple popular tools into a single, efficient command-line interface.

## Features ✨

- **Multiple Tools**: Integrates subfinder, amass, findomain, assetfinder, and more
- **API Integration**: SecurityTrails, VirusTotal, and Censys APIs for enhanced discovery
- **HTTP & Port Scanning**: Built-in httpx and smap integration
- **Concurrent Processing**: Multi-threaded execution for faster results
- **Multiple Output Formats**: JSON, TXT, and beautiful HTML reports
- **Flexible Input**: Support for both single domain and multiple domains
- **Easy Configuration**: YAML config files with CLI override support

## Quick Start 🚀

### Install

```bash
go install github.com/itszeeshan/subdomainx@latest
```

### Basic Usage

**Single domain:**

```bash
subdomainx --subfinder --httpx example.com
```

**Multiple domains:**

```bash
echo "example.com" > domains.txt
subdomainx --wildcard domains.txt --format html
```

**With API integration:**

```bash
# Set API keys
export SECURITYTRAILS_API_KEY="your_key"
export VIRUSTOTAL_API_KEY="your_key"
export CENSYS_API_ID="your_id"
export CENSYS_SECRET="your_secret"

# Use APIs
subdomainx --securitytrails --virustotal --censys example.com
```

> **Important**: Flags must be placed before the domain argument:
>
> ```bash
> subdomainx --tools domain.com  # Correct
> subdomainx domain.com --tools  # Incorrect
> ```

## Supported Tools 🛠️

### Enumeration Tools

subfinder, amass, findomain, assetfinder, sublist3r, knockpy, dnsrecon, fierce, massdns, altdns, waybackurls, linkheader

### API Tools

SecurityTrails, VirusTotal, Censys

### Scanning Tools

httpx (HTTP scanning), smap (port scanning)

## Tool Management

```bash
# Check available tools
subdomainx --check-tools

# Get installation help
subdomainx --install-tools
```

## Documentation 📚

For comprehensive documentation, examples, and advanced usage:

**🌐 [View Full Documentation](https://subdomainx.vercel.app)**

- [Installation Guide](https://subdomainx.vercel.app/installation)
- [CLI Reference](https://subdomainx.vercel.app/cli-reference)
- [Examples](https://subdomainx.vercel.app/examples)
- [Configuration](https://subdomainx.vercel.app/configuration)
- [Supported Tools](https://subdomainx.vercel.app/supported-tools)

## Examples 💡

```bash
# Single domain with specific tools
subdomainx --subfinder --amass --httpx example.com

# Multiple domains with HTML report
subdomainx --wildcard domains.txt --format html --name my_scan

# High-performance scan
subdomainx --threads 20 --timeout 60 example.com

# With filters
subdomainx --httpx --status-codes 200,301,302 --ports 80,443 example.com

# API-powered enumeration
subdomainx --securitytrails --virustotal --censys example.com

# Mixed approach (CLI tools + APIs)
subdomainx --subfinder --securitytrails --waybackurls --linkheader example.com
```

## Contributing 🤝

Contributions are welcome! Please feel free to submit a Pull Request.

## License 📄

MIT License - see [LICENSE](LICENSE) for details.

## Disclaimer ⚠️

This tool is for authorized security testing only. Always ensure proper authorization before scanning any domain.

---

**Happy Hunting! 🎯**
