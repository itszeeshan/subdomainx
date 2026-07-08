package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Tool represents a subdomain enumeration or scanning tool
type Tool struct {
	Name        string
	Command     string
	Description string
	InstallCmd  map[string]string // OS -> install command
	Required    bool
}

// GetRequiredTools returns a list of tools that SubdomainX can use
func GetRequiredTools() []Tool {
	return []Tool{
		{
			Name:        "subfinder",
			Command:     "subfinder",
			Description: "Fast subdomain discovery tool",
			InstallCmd: map[string]string{
				"linux":   "go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest",
				"darwin":  "go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest",
				"windows": "go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest",
			},
			Required: false,
		},
		{
			Name:        "amass",
			Command:     "amass",
			Description: "Comprehensive subdomain enumeration",
			InstallCmd: map[string]string{
				"linux":   "go install -v github.com/owasp-amass/amass/v4/...@master",
				"darwin":  "brew install amass || go install -v github.com/owasp-amass/amass/v4/...@master",
				"windows": "go install -v github.com/owasp-amass/amass/v4/...@master",
			},
			Required: false,
		},
		{
			Name:        "findomain",
			Command:     "findomain",
			Description: "Fast and cross-platform subdomain discovery",
			InstallCmd: map[string]string{
				"linux":   "wget https://github.com/Findomain/Findomain/releases/latest/download/findomain-linux -O /tmp/findomain && chmod +x /tmp/findomain && sudo mv /tmp/findomain /usr/local/bin/",
				"darwin":  "brew install findomain",
				"windows": "Download from: https://github.com/Findomain/Findomain/releases/latest",
			},
			Required: false,
		},
		{
			Name:        "assetfinder",
			Command:     "assetfinder",
			Description: "Find subdomains related to a domain",
			InstallCmd: map[string]string{
				"linux":   "go install github.com/tomnomnom/assetfinder@latest",
				"darwin":  "go install github.com/tomnomnom/assetfinder@latest",
				"windows": "go install github.com/tomnomnom/assetfinder@latest",
			},
			Required: false,
		},
		{
			Name:        "sublist3r",
			Command:     "sublist3r",
			Description: "Subdomain enumeration using OSINT",
			InstallCmd: map[string]string{
				"linux":   "pipx install sublist3r 2>/dev/null || pip3 install --user sublist3r",
				"darwin":  "pipx install sublist3r 2>/dev/null || pip3 install --user sublist3r",
				"windows": "pip install sublist3r",
			},
			Required: false,
		},
		{
			Name:        "knockpy",
			Command:     "knockpy",
			Description: "Subdomain enumeration tool",
			InstallCmd: map[string]string{
				"linux":   "pipx install knockpy 2>/dev/null || pip3 install --user knockpy",
				"darwin":  "pipx install knockpy 2>/dev/null || pip3 install --user knockpy",
				"windows": "pip install knockpy",
			},
			Required: false,
		},
		{
			Name:        "dnsrecon",
			Command:     "dnsrecon",
			Description: "DNS enumeration and reconnaissance",
			InstallCmd: map[string]string{
				"linux":   "pipx install dnsrecon 2>/dev/null || pip3 install --user dnsrecon",
				"darwin":  "pipx install dnsrecon 2>/dev/null || pip3 install --user dnsrecon",
				"windows": "pip install dnsrecon",
			},
			Required: false,
		},
		{
			Name:        "fierce",
			Command:     "fierce",
			Description: "DNS reconnaissance tool",
			InstallCmd: map[string]string{
				"linux":   "pipx install fierce 2>/dev/null || pip3 install --user fierce",
				"darwin":  "pipx install fierce 2>/dev/null || pip3 install --user fierce",
				"windows": "pip install fierce",
			},
			Required: false,
		},
		{
			Name:        "massdns",
			Command:     "massdns",
			Description: "High-performance DNS stub resolver",
			InstallCmd: map[string]string{
				"linux":   "git clone https://github.com/blechschmidt/massdns.git /tmp/massdns && cd /tmp/massdns && make && sudo cp bin/massdns /usr/local/bin/",
				"darwin":  "brew install massdns",
				"windows": "Download and compile from: https://github.com/blechschmidt/massdns",
			},
			Required: false,
		},
		{
			Name:        "altdns",
			Command:     "altdns",
			Description: "Subdomain permutation and alteration",
			InstallCmd: map[string]string{
				"linux":   "pipx install py-altdns 2>/dev/null || pip3 install --user py-altdns",
				"darwin":  "pipx install py-altdns 2>/dev/null || pip3 install --user py-altdns",
				"windows": "pip install py-altdns",
			},
			Required: false,
		},
		{
			Name:        "httpx",
			Command:     "httpx",
			Description: "Fast and multi-purpose HTTP probe",
			InstallCmd: map[string]string{
				"linux":   "go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest",
				"darwin":  "go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest",
				"windows": "go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest",
			},
			Required: false,
		},
		{
			Name:        "smap",
			Command:     "smap",
			Description: "Shodan-based port scanner",
			InstallCmd: map[string]string{
				"linux":   "go install -v github.com/s0md3v/smap/cmd/smap@latest",
				"darwin":  "go install -v github.com/s0md3v/smap/cmd/smap@latest",
				"windows": "go install -v github.com/s0md3v/smap/cmd/smap@latest",
			},
			Required: false,
		},
		{
			Name:        "nmap",
			Command:     "nmap",
			Description: "Network exploration and port scanning",
			InstallCmd: map[string]string{
				"linux":   "sudo apt-get install -y nmap || sudo yum install -y nmap",
				"darwin":  "brew install nmap",
				"windows": "Download from: https://nmap.org/download.html",
			},
			Required: false,
		},
		{
			Name:        "waybackurls",
			Command:     "waybackurls",
			Description: "Fetch URLs from Wayback Machine for subdomain discovery",
			InstallCmd: map[string]string{
				"linux":   "go install github.com/tomnomnom/waybackurls@latest",
				"darwin":  "go install github.com/tomnomnom/waybackurls@latest",
				"windows": "go install github.com/tomnomnom/waybackurls@latest",
			},
			Required: false,
		},
		// --- API-based tools (no binary, need env vars) ---
		{
			Name:        "securitytrails",
			Command:     "securitytrails",
			Description: "SecurityTrails API for subdomain enumeration",
			InstallCmd: map[string]string{
				"linux":   "Set SECURITYTRAILS_API_KEY environment variable",
				"darwin":  "Set SECURITYTRAILS_API_KEY environment variable",
				"windows": "Set SECURITYTRAILS_API_KEY environment variable",
			},
			Required: false,
		},
		{
			Name:        "virustotal",
			Command:     "virustotal",
			Description: "VirusTotal API for subdomain enumeration",
			InstallCmd: map[string]string{
				"linux":   "Set VIRUSTOTAL_API_KEY environment variable",
				"darwin":  "Set VIRUSTOTAL_API_KEY environment variable",
				"windows": "Set VIRUSTOTAL_API_KEY environment variable",
			},
			Required: false,
		},
		{
			Name:        "censys",
			Command:     "censys",
			Description: "Censys API for subdomain enumeration",
			InstallCmd: map[string]string{
				"linux":   "Set CENSYS_API_ID and CENSYS_SECRET environment variables",
				"darwin":  "Set CENSYS_API_ID and CENSYS_SECRET environment variables",
				"windows": "Set CENSYS_API_ID and CENSYS_SECRET environment variables",
			},
			Required: false,
		},
		// --- Built-in / public API tools (always available) ---
		{
			Name:        "linkheader",
			Command:     "linkheader",
			Description: "Discover subdomains from HTTP Link headers",
			InstallCmd: map[string]string{
				"linux":   "Built-in (no installation required)",
				"darwin":  "Built-in (no installation required)",
				"windows": "Built-in (no installation required)",
			},
			Required: false,
		},
		{
			Name:        "crtsh",
			Command:     "crtsh",
			Description: "Certificate Transparency database for subdomain discovery",
			InstallCmd: map[string]string{
				"linux":   "Built-in (no installation required)",
				"darwin":  "Built-in (no installation required)",
				"windows": "Built-in (no installation required)",
			},
			Required: false,
		},
		{
			Name:        "urlscan",
			Command:     "urlscan",
			Description: "URLScan.io API for subdomain enumeration",
			InstallCmd: map[string]string{
				"linux":   "Set URLSCAN_API_KEY environment variable (optional)",
				"darwin":  "Set URLSCAN_API_KEY environment variable (optional)",
				"windows": "Set URLSCAN_API_KEY environment variable (optional)",
			},
			Required: false,
		},
		{
			Name:        "hackertarget",
			Command:     "hackertarget",
			Description: "HackerTarget API for subdomain enumeration",
			InstallCmd: map[string]string{
				"linux":   "Set HACKERTARGET_API_KEY environment variable (optional)",
				"darwin":  "Set HACKERTARGET_API_KEY environment variable (optional)",
				"windows": "Set HACKERTARGET_API_KEY environment variable (optional)",
			},
			Required: false,
		},
	}
}

// CheckToolAvailability checks if a tool is available in PATH or as an API
func CheckToolAvailability(toolName string) bool {
	switch toolName {
	case "securitytrails":
		return strings.TrimSpace(os.Getenv("SECURITYTRAILS_API_KEY")) != ""
	case "virustotal":
		return strings.TrimSpace(os.Getenv("VIRUSTOTAL_API_KEY")) != ""
	case "censys":
		apiID := strings.TrimSpace(os.Getenv("CENSYS_API_ID"))
		secret := strings.TrimSpace(os.Getenv("CENSYS_SECRET"))
		return apiID != "" && secret != ""
	case "linkheader", "crtsh", "urlscan", "hackertarget":
		return true
	default:
		_, err := exec.LookPath(toolName)
		return err == nil
	}
}

// CheckAllTools checks the availability of all tools and returns missing ones
func CheckAllTools() ([]Tool, []Tool) {
	tools := GetRequiredTools()
	var available, missing []Tool

	for _, tool := range tools {
		if CheckToolAvailability(tool.Command) {
			available = append(available, tool)
		} else {
			missing = append(missing, tool)
		}
	}

	return available, missing
}

// prerequisite describes a required program for installing tools.
type prerequisite struct {
	Name    string
	Command string
	Version string // populated after check
}

// checkPrerequisites detects which package managers / runtimes are available.
func checkPrerequisites() []prerequisite {
	prereqs := []prerequisite{
		{Name: "go", Command: "go"},
		{Name: "pip3", Command: "pip3"},
		{Name: "git", Command: "git"},
	}
	if runtime.GOOS == "darwin" {
		prereqs = append(prereqs, prerequisite{Name: "brew", Command: "brew"})
	}

	for i := range prereqs {
		path, err := exec.LookPath(prereqs[i].Command)
		if err != nil {
			continue
		}
		_ = path
		out, err := exec.Command(prereqs[i].Command, "--version").CombinedOutput()
		if err == nil {
			line := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0]
			prereqs[i].Version = line
		} else {
			prereqs[i].Version = "installed"
		}
	}
	return prereqs
}

// isBuiltinOrAPI returns true if the tool's install command is not a real
// installable command (i.e. it's built-in or requires only an env var).
func isBuiltinOrAPI(tool Tool) bool {
	cmd, ok := tool.InstallCmd[runtime.GOOS]
	if !ok {
		return true
	}
	return strings.HasPrefix(cmd, "Built-in") || strings.HasPrefix(cmd, "Set ")
}

// isAPITool returns true if the tool only needs an API key, not a binary.
func isAPITool(tool Tool) bool {
	cmd, ok := tool.InstallCmd[runtime.GOOS]
	if !ok {
		return false
	}
	return strings.HasPrefix(cmd, "Set ")
}

// isDownloadOnly returns true if the install command is just a download link.
func isDownloadOnly(tool Tool) bool {
	cmd, ok := tool.InstallCmd[runtime.GOOS]
	if !ok {
		return true
	}
	return strings.HasPrefix(cmd, "Download")
}

// needsPrereq returns the prerequisite command needed to install a tool.
func needsPrereq(tool Tool) string {
	cmd, ok := tool.InstallCmd[runtime.GOOS]
	if !ok {
		return ""
	}
	if strings.HasPrefix(cmd, "go install") {
		return "go"
	}
	if strings.HasPrefix(cmd, "pipx ") || strings.HasPrefix(cmd, "pip3 ") || strings.HasPrefix(cmd, "pip ") {
		return "pip3"
	}
	if strings.HasPrefix(cmd, "brew ") {
		return "brew"
	}
	if strings.Contains(cmd, "git clone") {
		return "git"
	}
	return "" // wget/curl/sudo — no special prereq check
}

// InstallTools checks for missing tools and installs them automatically.
// This is the main entry point for `--install-tools`.
func InstallTools() error {
	fmt.Println("\n🔧 SubdomainX Tool Installer")
	fmt.Println("=============================")
	fmt.Println()

	// --- Prerequisites ---
	fmt.Println("Checking prerequisites...")
	prereqs := checkPrerequisites()
	prereqAvail := make(map[string]bool)
	for _, p := range prereqs {
		if p.Version != "" {
			prereqAvail[p.Name] = true
			fmt.Printf("  ✅ %s (%s)\n", p.Name, truncate(p.Version, 50))
		} else {
			prereqAvail[p.Name] = false
			fmt.Printf("  ❌ %s (not found)\n", p.Name)
		}
	}
	fmt.Println()

	// --- Categorize tools ---
	available, missing := CheckAllTools()

	var installable []Tool
	var apiTools []Tool
	var downloadOnly []Tool
	var missingPrereq []Tool

	for _, tool := range missing {
		if isBuiltinOrAPI(tool) {
			if isAPITool(tool) {
				apiTools = append(apiTools, tool)
			}
			continue
		}
		if isDownloadOnly(tool) {
			downloadOnly = append(downloadOnly, tool)
			continue
		}
		req := needsPrereq(tool)
		if req != "" && !prereqAvail[req] {
			missingPrereq = append(missingPrereq, tool)
			continue
		}
		installable = append(installable, tool)
	}

	fmt.Println("Checking tool status...")
	fmt.Printf("  ✅ %d tools already available\n", len(available))
	if len(installable) > 0 {
		fmt.Printf("  📦 %d tools to install\n", len(installable))
	}
	if len(apiTools) > 0 {
		fmt.Printf("  🔑 %d API tools need configuration\n", len(apiTools))
	}
	if len(missingPrereq) > 0 {
		fmt.Printf("  ⚠️  %d tools skipped (missing prerequisite)\n", len(missingPrereq))
	}
	if len(downloadOnly) > 0 {
		fmt.Printf("  ⬇️  %d tools need manual download\n", len(downloadOnly))
	}
	fmt.Println()

	if len(installable) == 0 && len(apiTools) == 0 && len(downloadOnly) == 0 && len(missingPrereq) == 0 {
		fmt.Println("✅ All tools are already installed!")
		return nil
	}

	// --- Install ---
	installed := 0
	failed := 0

	if len(installable) > 0 {
		fmt.Println("Installing tools...")
		for i, tool := range installable {
			installCmd := tool.InstallCmd[runtime.GOOS]
			fmt.Printf("  [%d/%d] %s ", i+1, len(installable), tool.Name)

			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/C", installCmd)
			} else {
				cmd = exec.Command("sh", "-c", installCmd)
			}

			// For sudo commands, connect stdin so password prompt works
			if strings.Contains(installCmd, "sudo") {
				cmd.Stdin = os.Stdin
			}

			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("❌ failed\n")
				// Show first line of error for debugging
				errLine := strings.SplitN(strings.TrimSpace(string(output)), "\n", 2)[0]
				if errLine != "" {
					fmt.Printf("         %s\n", truncate(errLine, 80))
				}
				failed++
				continue
			}

			// Verify it's actually available now
			if CheckToolAvailability(tool.Command) {
				fmt.Printf("✅ installed\n")
				installed++
			} else {
				fmt.Printf("⚠️  completed but not found in PATH\n")
				failed++
			}
		}
		fmt.Println()
	}

	// --- Skipped tools (missing prereq) ---
	if len(missingPrereq) > 0 {
		fmt.Println("Skipped (install the prerequisite first):")
		for _, tool := range missingPrereq {
			req := needsPrereq(tool)
			fmt.Printf("  • %s — needs %s\n", tool.Name, req)
		}
		fmt.Println()
	}

	// --- Download-only tools ---
	if len(downloadOnly) > 0 {
		fmt.Println("Manual download required:")
		for _, tool := range downloadOnly {
			cmd := tool.InstallCmd[runtime.GOOS]
			fmt.Printf("  • %s — %s\n", tool.Name, cmd)
		}
		fmt.Println()
	}

	// --- API keys ---
	if len(apiTools) > 0 {
		fmt.Println("API Keys (set these environment variables):")
		apiInfo := map[string]string{
			"securitytrails": "SECURITYTRAILS_API_KEY    — https://securitytrails.com/",
			"virustotal":     "VIRUSTOTAL_API_KEY        — https://virustotal.com/",
			"censys":         "CENSYS_API_ID + CENSYS_SECRET — https://search.censys.io/",
		}
		for _, tool := range apiTools {
			if info, ok := apiInfo[tool.Name]; ok {
				fmt.Printf("  • %s\n", info)
			}
		}
		fmt.Println()
	}

	// --- Summary ---
	finalAvail, finalMissing := CheckAllTools()
	fmt.Printf("Summary: %d/%d tools ready", len(finalAvail), len(finalAvail)+len(finalMissing))
	if len(apiTools) > 0 {
		fmt.Printf(", %d need API keys", len(apiTools))
	}
	if failed > 0 {
		fmt.Printf(", %d failed", failed)
	}
	fmt.Printf("\n")

	if installed > 0 {
		fmt.Printf("\n💡 Run 'subdomainx --check-tools' to verify.\n")
	}

	return nil
}

// DisplayToolStatus displays the status of all tools
func DisplayToolStatus() {
	available, missing := CheckAllTools()

	fmt.Println("\n🔧 Tool Status:")
	fmt.Println("================")

	if len(available) > 0 {
		fmt.Printf("\n✅ Available tools (%d):\n", len(available))
		for _, tool := range available {
			fmt.Printf("  • %s - %s\n", tool.Name, tool.Description)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("\n❌ Missing tools (%d):\n", len(missing))
		for _, tool := range missing {
			fmt.Printf("  • %s - %s\n", tool.Name, tool.Description)
		}
		fmt.Println("\n💡 Run with --install-tools to install them automatically")
	}

	fmt.Printf("\nTotal: %d available, %d missing\n", len(available), len(missing))
}

// truncate shortens s to max characters, appending "..." if truncated.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
