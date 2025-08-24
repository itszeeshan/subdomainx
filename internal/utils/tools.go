package utils

import (
	"bufio"
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
				"darwin":  "brew install findomain || (wget https://github.com/Findomain/Findomain/releases/latest/download/findomain-osx -O /tmp/findomain && chmod +x /tmp/findomain && sudo mv /tmp/findomain /usr/local/bin/)",
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
				"linux":   "pip3 install sublist3r || (git clone https://github.com/aboul3la/Sublist3r.git /tmp/Sublist3r && cd /tmp/Sublist3r && pip3 install -r requirements.txt && sudo cp sublist3r.py /usr/local/bin/sublist3r && sudo chmod +x /usr/local/bin/sublist3r)",
				"darwin":  "pip3 install sublist3r || (git clone https://github.com/aboul3la/Sublist3r.git /tmp/Sublist3r && cd /tmp/Sublist3r && pip3 install -r requirements.txt && sudo cp sublist3r.py /usr/local/bin/sublist3r && sudo chmod +x /usr/local/bin/sublist3r)",
				"windows": "pip install sublist3r",
			},
			Required: false,
		},
		{
			Name:        "knockpy",
			Command:     "knockpy",
			Description: "Subdomain enumeration tool",
			InstallCmd: map[string]string{
				"linux":   "pip3 install knockpy",
				"darwin":  "pip3 install knockpy",
				"windows": "pip install knockpy",
			},
			Required: false,
		},
		{
			Name:        "dnsrecon",
			Command:     "dnsrecon",
			Description: "DNS enumeration and reconnaissance",
			InstallCmd: map[string]string{
				"linux":   "pip3 install dnsrecon || sudo apt-get install dnsrecon",
				"darwin":  "pip3 install dnsrecon",
				"windows": "pip install dnsrecon",
			},
			Required: false,
		},
		{
			Name:        "fierce",
			Command:     "fierce",
			Description: "DNS reconnaissance tool",
			InstallCmd: map[string]string{
				"linux":   "pip3 install fierce",
				"darwin":  "pip3 install fierce",
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
				"darwin":  "brew install massdns || (git clone https://github.com/blechschmidt/massdns.git /tmp/massdns && cd /tmp/massdns && make && sudo cp bin/massdns /usr/local/bin/)",
				"windows": "Download and compile from: https://github.com/blechschmidt/massdns",
			},
			Required: false,
		},
		{
			Name:        "altdns",
			Command:     "altdns",
			Description: "Subdomain permutation and alteration",
			InstallCmd: map[string]string{
				"linux":   "pip3 install py-altdns",
				"darwin":  "pip3 install py-altdns",
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
			Name:        "nmap",
			Command:     "nmap",
			Description: "Network exploration and port scanning",
			InstallCmd: map[string]string{
				"linux":   "sudo apt-get install nmap || sudo yum install nmap",
				"darwin":  "brew install nmap",
				"windows": "Download from: https://nmap.org/download.html",
			},
			Required: false,
		},
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
			Name:        "threatcrowd",
			Command:     "threatcrowd",
			Description: "ThreatCrowd API for subdomain enumeration",
			InstallCmd: map[string]string{
				"linux":   "Built-in (no installation required)",
				"darwin":  "Built-in (no installation required)",
				"windows": "Built-in (no installation required)",
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
	// Special handling for API-based tools
	switch toolName {
	case "securitytrails":
		// Check if API key is configured
		apiKey := strings.TrimSpace(os.Getenv("SECURITYTRAILS_API_KEY"))
		return apiKey != ""
	case "virustotal":
		// Check if API key is configured
		apiKey := strings.TrimSpace(os.Getenv("VIRUSTOTAL_API_KEY"))
		return apiKey != ""
	case "censys":
		// Check if API credentials are configured
		apiID := strings.TrimSpace(os.Getenv("CENSYS_API_ID"))
		secret := strings.TrimSpace(os.Getenv("CENSYS_SECRET"))
		return apiID != "" && secret != ""
	case "linkheader":
		// Link header enumerator is built-in, always available
		return true
	case "crtsh":
		// crt.sh is a public API, always available
		return true
	case "urlscan":
		// URLScan.io can work without API key, but better with one
		_ = strings.TrimSpace(os.Getenv("URLSCAN_API_KEY"))
		return true // Always available, API key is optional
	case "threatcrowd":
		// ThreatCrowd is a public API, always available
		return true
	case "hackertarget":
		// HackerTarget can work without API key, but better with one
		_ = strings.TrimSpace(os.Getenv("HACKERTARGET_API_KEY"))
		return true // Always available, API key is optional
	default:
		// For command-line tools, check if they're in PATH
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

// PromptToolInstallation prompts the user to install missing tools
func PromptToolInstallation(missingTools []Tool) error {
	if len(missingTools) == 0 {
		fmt.Println("✅ All tools are available!")
		return nil
	}

	fmt.Printf("\n⚠️  Found %d missing tools:\n\n", len(missingTools))

	for i, tool := range missingTools {
		fmt.Printf("%d. %s - %s\n", i+1, tool.Name, tool.Description)
	}

	fmt.Print("\nWould you like to see installation instructions? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("\n⚠️  Some tools are missing. SubdomainX will skip unavailable tools during enumeration.")
		return nil
	}

	// Show installation instructions
	fmt.Printf("\n📦 Installation instructions for %s:\n\n", runtime.GOOS)

	for i, tool := range missingTools {
		fmt.Printf("--- %d. %s ---\n", i+1, tool.Name)

		if installCmd, exists := tool.InstallCmd[runtime.GOOS]; exists {
			fmt.Printf("Command: %s\n", installCmd)
		} else {
			fmt.Printf("Please visit the tool's official repository for installation instructions.\n")
		}
		fmt.Println()
	}

	fmt.Print("Would you like to automatically install tools that support auto-installation? (y/n): ")
	autoResponse, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	autoResponse = strings.TrimSpace(strings.ToLower(autoResponse))
	if autoResponse == "y" || autoResponse == "yes" {
		return autoInstallTools(missingTools)
	}

	fmt.Println("\n💡 After installing the tools, run SubdomainX again.")
	fmt.Println("💡 You can also disable specific tools in the config file if you don't want to install them.")

	return nil
}

// autoInstallTools attempts to automatically install tools
func autoInstallTools(tools []Tool) error {
	fmt.Println("\n🔧 Starting automatic installation...")

	for _, tool := range tools {
		installCmd, exists := tool.InstallCmd[runtime.GOOS]
		if !exists {
			fmt.Printf("⏭️  Skipping %s (manual installation required)\n", tool.Name)
			continue
		}

		// Skip tools that require manual steps or sudo
		if strings.Contains(installCmd, "sudo") || strings.Contains(installCmd, "Download") {
			fmt.Printf("⏭️  Skipping %s (requires manual installation)\n", tool.Name)
			continue
		}

		fmt.Printf("📦 Installing %s...\n", tool.Name)

		// Execute installation command
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", installCmd)
		} else {
			cmd = exec.Command("sh", "-c", installCmd)
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Failed to install %s: %v\n", tool.Name, err)
			continue
		}

		// Verify installation
		if CheckToolAvailability(tool.Command) {
			fmt.Printf("✅ Successfully installed %s\n", tool.Name)
		} else {
			fmt.Printf("⚠️  %s installation completed but tool not found in PATH\n", tool.Name)
		}
		fmt.Println()
	}

	fmt.Println("🎉 Automatic installation completed!")
	fmt.Println("💡 Some tools may require manual installation or PATH configuration.")

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
		fmt.Println("\n💡 Run with --install-tools to see installation instructions")
	}

	fmt.Printf("\nTotal: %d available, %d missing", len(available), len(missing))
}
