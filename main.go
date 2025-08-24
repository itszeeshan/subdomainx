package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/enumerator"
	"github.com/itszeeshan/subdomainx/internal/output"
	"github.com/itszeeshan/subdomainx/internal/scanner"
	"github.com/itszeeshan/subdomainx/internal/types"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

func main() {
	// Parse command line flags
	var (
		showVersion     = flag.Bool("version", false, "Show version information")
		showHelp        = flag.Bool("help", false, "Show help information")
		checkTools      = flag.Bool("check-tools", false, "Check tool availability")
		installTools    = flag.Bool("install-tools", false, "Show tool installation instructions")
		configFile      = flag.String("config", "", "Path to configuration file (optional)")
		wildcardFile    = flag.String("wildcard", "", "Path to wildcard file containing domains")
		uniqueName      = flag.String("name", "scan", "Unique name for output files")
		outputFormat    = flag.String("format", "", "Output format (json, txt, html)")
		outputDir       = flag.String("output", "output", "Output directory")
		threads         = flag.Int("threads", 10, "Number of threads")
		retries         = flag.Int("retries", 3, "Number of retry attempts")
		timeout         = flag.Int("timeout", 30, "Timeout in seconds")
		rateLimit       = flag.Int("rate-limit", 100, "Rate limit per second")
		wordlist        = flag.String("wordlist", "", "Custom wordlist file for brute-forcing")
		resume          = flag.String("resume", "", "Resume scan from checkpoint (scan ID)")
		listCheckpoints = flag.Bool("list-checkpoints", false, "List available checkpoints")
		verbose         = flag.Bool("verbose", false, "Verbose output")

		// Filter flags
		statusCodes = flag.String("status-codes", "", "Filter by HTTP status codes (e.g., '200,301,302')")
		ports       = flag.String("ports", "", "Filter by ports (e.g., '80,443,8080')")

		// Tool-specific flags
		useSubfinder      = flag.Bool("subfinder", false, "Use subfinder tool")
		useAmass          = flag.Bool("amass", false, "Use amass tool")
		useFindomain      = flag.Bool("findomain", false, "Use findomain tool")
		useAssetfinder    = flag.Bool("assetfinder", false, "Use assetfinder tool")
		useSublist3r      = flag.Bool("sublist3r", false, "Use sublist3r tool")
		useKnockpy        = flag.Bool("knockpy", false, "Use knockpy tool")
		useDnsrecon       = flag.Bool("dnsrecon", false, "Use dnsrecon tool")
		useFierce         = flag.Bool("fierce", false, "Use fierce tool")
		useMassdns        = flag.Bool("massdns", false, "Use massdns tool")
		useAltdns         = flag.Bool("altdns", false, "Use altdns tool")
		useSecurityTrails = flag.Bool("securitytrails", false, "Use SecurityTrails API")
		useVirusTotal     = flag.Bool("virustotal", false, "Use VirusTotal API")
		useCensys         = flag.Bool("censys", false, "Use Censys API")
		useWaybackURLs    = flag.Bool("waybackurls", false, "Use waybackurls tool")
		useLinkHeader     = flag.Bool("linkheader", false, "Use Link Header enumeration")
		useHttpx          = flag.Bool("httpx", false, "Use httpx for HTTP scanning")
		useSmap           = flag.Bool("smap", false, "Use smap for port scanning")
	)
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Println("SubdomainX v1.2.0")
		fmt.Println("All-in-one subdomain enumeration tool")
		return
	}

	// Show help
	if *showHelp {
		showUsage()
		return
	}

	// Check tools (doesn't require domain or wildcard file)
	if *checkTools {
		utils.DisplayToolStatus()
		return
	}

	// Show installation instructions (doesn't require domain or wildcard file)
	if *installTools {
		_, missing := utils.CheckAllTools()
		if err := utils.PromptToolInstallation(missing); err != nil {
			log.Fatalf("Failed to show installation instructions: %v", err)
		}
		return
	}

	// List available checkpoints
	if *listCheckpoints {
		checkpoints, err := utils.ListCheckpoints(*outputDir)
		if err != nil {
			log.Fatalf("Failed to list checkpoints: %v", err)
		}

		if len(checkpoints) == 0 {
			fmt.Println("📋 No checkpoints found.")
		} else {
			fmt.Println("📋 Available checkpoints:")
			for _, cp := range checkpoints {
				fmt.Printf("  • %s\n", cp)
			}
			fmt.Printf("\n💡 Resume with: subdomainx --resume <scan_id>\n")
		}
		return
	}

	// Display banner (always show, but can be controlled with verbose)
	showBanner()

	// Start resource monitoring if verbose mode is enabled
	if *verbose {
		utils.StartResourceMonitoring()
		defer utils.StopResourceMonitoring()
	}

	// Check tool availability and prompt for installation
	available, missing := utils.CheckAllTools()
	if len(missing) > 0 && *verbose {
		fmt.Printf("⚠️  Warning: %d tools are missing. SubdomainX will work with available tools.\n", len(missing))
		fmt.Println("💡 Use --install-tools to see installation instructions.")
		fmt.Printf("✅ Available tools: %d\n\n", len(available))
	}

	// Build configuration from CLI flags and optional config file
	cfg := &config.Config{
		OutputDir:    *outputDir,
		OutputFormat: "json", // Default format
		Threads:      *threads,
		Retries:      *retries,
		Timeout:      *timeout,
		RateLimit:    *rateLimit,
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	// Check if we have a domain argument
	args := flag.Args()
	hasDomainArg := len(args) > 0

	// Override with CLI flags (CLI takes highest precedence)
	if *wildcardFile != "" {
		cfg.WildcardFile = *wildcardFile
	}
	if *uniqueName != "" {
		cfg.UniqueName = *uniqueName
	}
	if *outputFormat != "" {
		cfg.OutputFormat = *outputFormat
	}
	if *wordlist != "" {
		cfg.Wordlist = *wordlist
	}

	// Override filters with CLI flags
	if *statusCodes != "" {
		cfg.Filters["status_code"] = *statusCodes
	}
	if *ports != "" {
		cfg.Filters["ports"] = *ports
	}

	// Check if we have either a wildcard file, a single domain argument, or are resuming
	if cfg.WildcardFile == "" && len(args) == 0 && *resume == "" {
		log.Fatalf("Error: Either --wildcard file, a domain argument, or --resume is required. Use --help for usage information.")
	}

	// If no wildcard file but domain argument provided, create a temporary file
	if cfg.WildcardFile == "" && len(args) > 0 {
		// Use the first argument as the domain
		domain := args[0]

		// Validate the domain
		if err := utils.ValidateDomain(domain); err != nil {
			log.Fatalf("Error: Invalid domain '%s': %v", domain, err)
		}

		// Create a temporary file with the single domain
		tmpFile, err := os.CreateTemp("", "subdomainx_domain_*.txt")
		if err != nil {
			log.Fatalf("Error: Failed to create temporary file: %v", err)
		}
		defer os.Remove(tmpFile.Name()) // Clean up temp file

		// Write the domain to the temporary file
		if _, err := tmpFile.WriteString(domain + "\n"); err != nil {
			log.Fatalf("Error: Failed to write domain to temporary file: %v", err)
		}
		tmpFile.Close()

		// Set the wildcard file to our temporary file
		cfg.WildcardFile = tmpFile.Name()

		// Use the domain as the unique name if not specified
		if cfg.UniqueName == "scan" {
			cfg.UniqueName = domain
		}
	}

	// Handle tool selection
	// If any specific tools are specified via CLI, use only those
	// Otherwise, use all available tools from config
	specificToolsSelected := *useSubfinder || *useAmass || *useFindomain || *useAssetfinder ||
		*useSublist3r || *useKnockpy || *useDnsrecon || *useFierce || *useMassdns || *useAltdns ||
		*useSecurityTrails || *useVirusTotal || *useCensys || *useWaybackURLs || *useLinkHeader || *useHttpx || *useSmap

	// Load config file if specified (optional)
	if *configFile != "" {
		fileCfg, err := config.LoadConfigFromFile(*configFile)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
		// Merge config file with CLI flags (CLI takes precedence)
		cfg = mergeConfig(fileCfg, cfg)
	} else if !hasDomainArg && *resume == "" {
		// Only load default config if we don't have a domain argument and not resuming
		// This prevents loading default wildcard_file when using single domain mode or resuming
		if defaultCfg, err := config.LoadConfig(); err == nil {
			cfg = mergeConfig(defaultCfg, cfg)
		}
	}

	// Now handle tool selection - this MUST happen after config loading
	if specificToolsSelected {
		// Clear all tools and set only the specified ones
		// This completely overrides any config file settings
		cfg.Tools = make(map[string]bool)
		cfg.Tools["subfinder"] = *useSubfinder
		cfg.Tools["amass"] = *useAmass
		cfg.Tools["findomain"] = *useFindomain
		cfg.Tools["assetfinder"] = *useAssetfinder
		cfg.Tools["sublist3r"] = *useSublist3r
		cfg.Tools["knockpy"] = *useKnockpy
		cfg.Tools["dnsrecon"] = *useDnsrecon
		cfg.Tools["fierce"] = *useFierce
		cfg.Tools["massdns"] = *useMassdns
		cfg.Tools["altdns"] = *useAltdns
		cfg.Tools["securitytrails"] = *useSecurityTrails
		cfg.Tools["virustotal"] = *useVirusTotal
		cfg.Tools["censys"] = *useCensys
		cfg.Tools["waybackurls"] = *useWaybackURLs
		cfg.Tools["linkheader"] = *useLinkHeader
		cfg.Tools["httpx"] = *useHttpx
		cfg.Tools["smap"] = *useSmap

		// Debug: Print selected tools
		if *verbose {
			fmt.Printf("🔧 CLI Tools selected: ")
			var selectedTools []string
			for tool, enabled := range cfg.Tools {
				if enabled {
					selectedTools = append(selectedTools, tool)
				}
			}
			fmt.Printf("%s\n", strings.Join(selectedTools, ", "))
		}

	} else {
		// If no specific tools selected, ensure all tools are enabled
		// This ensures the tool works even without a config file
		if cfg.Tools == nil {
			cfg.Tools = make(map[string]bool)
		}
		cfg.Tools["subfinder"] = true
		cfg.Tools["amass"] = true
		cfg.Tools["findomain"] = true
		cfg.Tools["assetfinder"] = true
		cfg.Tools["sublist3r"] = true
		cfg.Tools["knockpy"] = true
		cfg.Tools["dnsrecon"] = true
		cfg.Tools["fierce"] = true
		cfg.Tools["massdns"] = true
		cfg.Tools["altdns"] = true
		cfg.Tools["securitytrails"] = true
		cfg.Tools["virustotal"] = true
		cfg.Tools["censys"] = true
		cfg.Tools["waybackurls"] = true
		cfg.Tools["linkheader"] = true
		cfg.Tools["httpx"] = true
		cfg.Tools["smap"] = true
	}

	// Validate input
	if err := validateCLIInput(cfg); err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	// Create output directory
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Initialize DNS cache
	dnsCache := cache.NewDNSCache()

	// Handle resume functionality
	var checkpoint *utils.Checkpoint
	var results []types.SubdomainResult
	var httpResults []types.HTTPResult
	var portResults []types.PortResult
	var err error

	if *resume != "" {
		// Load checkpoint
		checkpoint, err = utils.LoadCheckpoint(*resume, cfg.OutputDir)
		if err != nil {
			log.Fatalf("Failed to load checkpoint: %v", err)
		}

		fmt.Printf("🔄 Resuming scan from checkpoint: %s\n", *resume)
		fmt.Printf("📊 Previous progress: %d/%d tasks completed\n",
			checkpoint.Progress.CompletedTasks, checkpoint.Progress.TotalTasks)
		fmt.Printf("🔍 Previous subdomains found: %d\n", len(checkpoint.Subdomains))
		fmt.Printf("🌐 Previous HTTP results: %d\n", len(checkpoint.HTTPResults))
		fmt.Printf("🔌 Previous port results: %d\n", len(checkpoint.PortResults))

		// Use results from checkpoint
		results = checkpoint.Subdomains
		httpResults = checkpoint.HTTPResults
		portResults = checkpoint.PortResults

		// Update config with checkpoint data
		if checkpoint.Domain != "" {
			cfg.WildcardFile = checkpoint.WildcardFile
		}
	} else {
		// Create new checkpoint
		scanID := cfg.UniqueName
		domain := ""
		if len(args) > 0 {
			domain = args[0]
			if scanID == "scan" {
				scanID = args[0] // Use domain as scan ID
			}
		}

		configMap := map[string]interface{}{
			"threads":   cfg.Threads,
			"retries":   cfg.Retries,
			"timeout":   cfg.Timeout,
			"rateLimit": cfg.RateLimit,
			"wordlist":  cfg.Wordlist,
		}

		checkpoint = utils.CreateCheckpoint(scanID, domain, cfg.WildcardFile, configMap)
	}

	// Initialize signal handler for graceful interruption
	signalHandler := utils.NewSignalHandler(checkpoint, cfg.OutputDir)
	signalHandler.Start()

	// Run subdomain enumeration (only if not resuming or if resuming but enumeration not complete)
	if *resume == "" || len(results) == 0 {
		results, err = enumerator.Run(cfg, dnsCache)
		if err != nil {
			checkpoint.MarkError(fmt.Sprintf("Enumeration failed: %v", err))
			if saveErr := utils.SaveCheckpoint(checkpoint, cfg.OutputDir); saveErr != nil {
				log.Printf("Warning: Failed to save checkpoint: %v", saveErr)
			}
			log.Fatalf("Enumeration failed: %v", err)
		}

		// Update checkpoint with enumeration results
		checkpoint.AddSubdomains(results)
		checkpoint.UpdateProgress(len(results), len(results))
		if saveErr := utils.SaveCheckpoint(checkpoint, cfg.OutputDir); saveErr != nil {
			log.Printf("Warning: Failed to save checkpoint: %v", saveErr)
		}
	}
	if err != nil {
		log.Fatalf("Enumeration failed: %v", err)
	}

	// Run HTTP scanning only if httpx is enabled
	if cfg.Tools["httpx"] && (*resume == "" || len(httpResults) == 0) {
		log.Println("🔍 Running HTTP scanning with httpx...")
		httpResults, err = scanner.RunHTTPx(cfg, results)
		if err != nil {
			log.Printf("HTTP scanning failed: %v", err)
		} else {
			log.Printf("✅ HTTP scanning completed: %d results", len(httpResults))
			// Update checkpoint with HTTP results
			checkpoint.AddHTTPResults(httpResults)
			if saveErr := utils.SaveCheckpoint(checkpoint, cfg.OutputDir); saveErr != nil {
				log.Printf("Warning: Failed to save checkpoint: %v", saveErr)
			}
		}
	}

	// Run port scanning only if smap is enabled
	if cfg.Tools["smap"] && (*resume == "" || len(portResults) == 0) {
		log.Println("🔍 Running port scanning with smap...")
		portResults, err = scanner.RunSmap(cfg, results)
		if err != nil {
			log.Printf("Port scanning failed: %v", err)
		} else {
			log.Printf("✅ Port scanning completed: %d results", len(portResults))
			// Update checkpoint with port results
			checkpoint.AddPortResults(portResults)
			if saveErr := utils.SaveCheckpoint(checkpoint, cfg.OutputDir); saveErr != nil {
				log.Printf("Warning: Failed to save checkpoint: %v", saveErr)
			}
		}
	}

	// Generate output
	if err := output.Generate(cfg, results, httpResults, portResults); err != nil {
		log.Fatalf("Failed to generate output: %v", err)
	}

	// Mark checkpoint as completed
	checkpoint.MarkCompleted()
	if saveErr := utils.SaveCheckpoint(checkpoint, cfg.OutputDir); saveErr != nil {
		log.Printf("Warning: Failed to save checkpoint: %v", saveErr)
	}

	log.Printf("Scan completed. Results saved to %s", cfg.OutputDir)
}

// showBanner displays the application banner
func showBanner() {
	fmt.Print(`
  ____        _         _                       _       __  __
 / ___| _   _| |__   __| | ___  _ __ ___   __ _(_)_ __  \ \/ /
 \___ \| | | | '_ \ / _' |/ _ \| '_ ' _ \ / _' | | '_ \  \  / 
  ___) | |_| | |_) | (_| | (_) | | | | | | (_| | | | | | /  \ 
 |____/ \__,_|_.__/ \__,_|\___/|_| |_| |_|\__,_|_|_| |_|/_/\_\

        🔍 All-in-one Subdomain Enumeration Tool
        📧 https://github.com/itszeeshan/subdomainx
`)
}

// showUsage displays usage information
func showUsage() {
	fmt.Println(`SubdomainX - All-in-one Subdomain Enumeration Tool

USAGE:
    subdomainx <domain> [OPTIONS]                    # Single domain scan
    subdomainx --wildcard <domains_file> [OPTIONS]   # Multiple domains scan

REQUIRED (choose one):
    <domain>              Target domain for single domain scan
    --wildcard FILE       Path to file containing target domains (one per line)

OPTIONS:
    --version              Show version information
    --help                 Show this help message
    --check-tools          Check availability of enumeration tools
    --install-tools        Show installation instructions for missing tools
    
    # Output Options
    --name NAME            Unique name for output files (default: scan)
    --format FORMAT        Output format: json, txt, html, zap, burp, nessus, csv (default: json)
    --output DIR           Output directory (default: output)
    
    # Performance Options
    --threads N            Number of concurrent threads (default: 10)
    --retries N            Number of retry attempts (default: 3)
    --timeout N            Timeout in seconds (default: 30)
    --rate-limit N         Rate limit per second (default: 100)
    --wordlist FILE        Custom wordlist file for brute-forcing
    --resume SCAN_ID       Resume scan from checkpoint (scan ID)
    --list-checkpoints     List available checkpoints
    
    # Filter Options
    --status-codes CODES   Filter by HTTP status codes (e.g., '200,301,302')
    --ports PORTS          Filter by ports (e.g., '80,443,8080')
    
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
    --securitytrails       Use SecurityTrails API
    --virustotal           Use VirusTotal API
    --censys               Use Censys API
    --waybackurls          Use waybackurls tool
    --linkheader           Use Link Header enumeration
    --httpx                Use httpx for HTTP scanning
    --smap                 Use smap for port scanning
    
    # Configuration
    --config FILE          Use custom configuration file (optional)
    --verbose              Enable verbose output

EXAMPLES:
    # Single domain scan with all available tools
    subdomainx example.com

    # Single domain with specific tools
    subdomainx example.com --subfinder --amass --httpx

    # Multiple domains scan
    subdomainx --wildcard domains.txt

    # Multiple domains with specific tools
    subdomainx --wildcard domains.txt --amass --subfinder --httpx

    # Generate HTML report for single domain
    subdomainx example.com --format html --name my_scan

    # High-performance scan
    subdomainx example.com --threads 20 --timeout 60

    # Custom wordlist scan
    subdomainx --wordlist /path/to/wordlist.txt example.com

    # Resume interrupted scan
    subdomainx --resume my_scan

    # List available checkpoints
    subdomainx --list-checkpoints

    # Check tool availability
    subdomainx --check-tools

    # Get installation help
    subdomainx --install-tools

CONFIGURATION:
    The YAML config file is optional. All settings can be specified via CLI flags.
    CLI flags take precedence over config file settings.

For more information, visit: https://github.com/itszeeshan/subdomainx`)
}

// mergeConfig merges two configs, with cfg2 taking precedence over cfg1
func mergeConfig(cfg1, cfg2 *config.Config) *config.Config {
	result := &config.Config{
		WildcardFile: cfg1.WildcardFile,
		UniqueName:   cfg1.UniqueName,
		OutputDir:    cfg1.OutputDir,
		OutputFormat: cfg1.OutputFormat,
		Threads:      cfg1.Threads,
		Retries:      cfg1.Retries,
		Timeout:      cfg1.Timeout,
		RateLimit:    cfg1.RateLimit,
		Wordlist:     cfg1.Wordlist,
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	// Copy tools from cfg1
	for k, v := range cfg1.Tools {
		result.Tools[k] = v
	}

	// Copy filters from cfg1
	for k, v := range cfg1.Filters {
		result.Filters[k] = v
	}

	// Override with cfg2 values (cfg2 takes precedence)
	if cfg2.WildcardFile != "" {
		result.WildcardFile = cfg2.WildcardFile
	}
	if cfg2.UniqueName != "" {
		result.UniqueName = cfg2.UniqueName
	}
	if cfg2.OutputDir != "" {
		result.OutputDir = cfg2.OutputDir
	}
	if cfg2.OutputFormat != "" {
		result.OutputFormat = cfg2.OutputFormat
	}
	if cfg2.Threads > 0 {
		result.Threads = cfg2.Threads
	}
	if cfg2.Retries > 0 {
		result.Retries = cfg2.Retries
	}
	if cfg2.Timeout > 0 {
		result.Timeout = cfg2.Timeout
	}
	if cfg2.RateLimit > 0 {
		result.RateLimit = cfg2.RateLimit
	}
	if cfg2.Wordlist != "" {
		result.Wordlist = cfg2.Wordlist
	}

	// Override tools from cfg2
	for k, v := range cfg2.Tools {
		result.Tools[k] = v
	}

	// Override filters from cfg2
	for k, v := range cfg2.Filters {
		result.Filters[k] = v
	}

	return result
}

// validateCLIInput validates the CLI configuration
func validateCLIInput(cfg *config.Config) error {
	// Validate wildcard file exists (skip for temporary files created for single domain and resume mode)
	if cfg.WildcardFile != "" && !strings.Contains(cfg.WildcardFile, "subdomainx_domain_") && !utils.FileExists(cfg.WildcardFile) {
		return fmt.Errorf("wildcard file not found: %s", cfg.WildcardFile)
	}

	// Validate output directory can be created
	if err := utils.EnsureDirectory(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Validate output format
	validFormats := map[string]bool{
		"json":   true,
		"txt":    true,
		"html":   true,
		"zap":    true,
		"burp":   true,
		"nessus": true,
		"csv":    true,
	}
	if !validFormats[cfg.OutputFormat] {
		return fmt.Errorf("invalid output format: %s. Supported formats: json, txt, html, zap, burp, nessus, csv", cfg.OutputFormat)
	}

	// Validate threads
	if cfg.Threads <= 0 {
		return fmt.Errorf("threads must be greater than 0")
	}

	// Validate retries
	if cfg.Retries < 0 {
		return fmt.Errorf("retries cannot be negative")
	}

	// Validate timeout
	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	// Validate rate limit
	if cfg.RateLimit <= 0 {
		return fmt.Errorf("rate limit must be greater than 0")
	}

	// Validate wordlist if specified
	if cfg.Wordlist != "" && !utils.FileExists(cfg.Wordlist) {
		return fmt.Errorf("wordlist file not found: %s", cfg.Wordlist)
	}

	return nil
}
