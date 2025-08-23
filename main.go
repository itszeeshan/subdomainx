package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/enumerator"
	"github.com/itszeeshan/subdomainx/internal/output"
	"github.com/itszeeshan/subdomainx/internal/scanner"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

func main() {
	// Parse command line flags
	var (
		showVersion  = flag.Bool("version", false, "Show version information")
		showHelp     = flag.Bool("help", false, "Show help information")
		checkTools   = flag.Bool("check-tools", false, "Check tool availability")
		installTools = flag.Bool("install-tools", false, "Show tool installation instructions")
		configFile   = flag.String("config", "", "Path to configuration file (optional)")
		wildcardFile = flag.String("wildcard", "", "Path to wildcard file containing domains")
		uniqueName   = flag.String("name", "scan", "Unique name for output files")
		outputFormat = flag.String("format", "", "Output format (json, txt, html)")
		outputDir    = flag.String("output", "output", "Output directory")
		threads      = flag.Int("threads", 10, "Number of threads")
		retries      = flag.Int("retries", 3, "Number of retry attempts")
		timeout      = flag.Int("timeout", 30, "Timeout in seconds")
		rateLimit    = flag.Int("rate-limit", 100, "Rate limit per second")
		verbose      = flag.Bool("v", false, "Verbose output")

		// Tool-specific flags
		useSubfinder   = flag.Bool("subfinder", false, "Use subfinder tool")
		useAmass       = flag.Bool("amass", false, "Use amass tool")
		useFindomain   = flag.Bool("findomain", false, "Use findomain tool")
		useAssetfinder = flag.Bool("assetfinder", false, "Use assetfinder tool")
		useSublist3r   = flag.Bool("sublist3r", false, "Use sublist3r tool")
		useKnockpy     = flag.Bool("knockpy", false, "Use knockpy tool")
		useDnsrecon    = flag.Bool("dnsrecon", false, "Use dnsrecon tool")
		useFierce      = flag.Bool("fierce", false, "Use fierce tool")
		useMassdns     = flag.Bool("massdns", false, "Use massdns tool")
		useAltdns      = flag.Bool("altdns", false, "Use altdns tool")
		useHttpx       = flag.Bool("httpx", false, "Use httpx for HTTP scanning")
		useSmap        = flag.Bool("smap", false, "Use smap for port scanning")
	)
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Println("SubdomainX v1.0.0")
		fmt.Println("All-in-one subdomain enumeration tool")
		return
	}

	// Show help
	if *showHelp {
		showUsage()
		return
	}

	// Check tools
	if *checkTools {
		utils.DisplayToolStatus()
		return
	}

	// Show installation instructions
	if *installTools {
		_, missing := utils.CheckAllTools()
		if err := utils.PromptToolInstallation(missing); err != nil {
			log.Fatalf("Failed to show installation instructions: %v", err)
		}
		return
	}

	// Display banner
	if *verbose {
		showBanner()
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

	// Load config file if specified (optional)
	if *configFile != "" {
		fileCfg, err := config.LoadConfigFromFile(*configFile)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
		// Merge config file with CLI flags (CLI takes precedence)
		cfg = mergeConfig(fileCfg, cfg)
	} else {
		// Try to load default config if no CLI config specified
		if defaultCfg, err := config.LoadConfig(); err == nil {
			cfg = mergeConfig(defaultCfg, cfg)
		}
	}

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

	// Validate that wildcard file is provided for scans
	if cfg.WildcardFile == "" {
		log.Fatalf("Error: --wildcard file is required for scanning. Use --help for usage information.")
	}

	// Handle tool selection
	// If any specific tools are specified via CLI, use only those
	// Otherwise, use all available tools from config
	specificToolsSelected := *useSubfinder || *useAmass || *useFindomain || *useAssetfinder ||
		*useSublist3r || *useKnockpy || *useDnsrecon || *useFierce || *useMassdns || *useAltdns

	if specificToolsSelected {
		// Clear all tools and set only the specified ones
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
		cfg.Tools["httpx"] = *useHttpx
		cfg.Tools["smap"] = *useSmap
	} else {
		// If no specific tools selected, enable all available tools by default
		// This ensures the tool works even without a config file
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

	// Run subdomain enumeration
	results, err := enumerator.Run(cfg, dnsCache)
	if err != nil {
		log.Fatalf("Enumeration failed: %v", err)
	}

	// Run HTTP scanning
	httpResults, err := scanner.RunHTTPx(cfg, results)
	if err != nil {
		log.Printf("HTTP scanning failed: %v", err)
	}

	// Run port scanning
	portResults, err := scanner.RunSmap(cfg, results)
	if err != nil {
		log.Printf("Port scanning failed: %v", err)
	}

	// Generate output
	if err := output.Generate(cfg, results, httpResults, portResults); err != nil {
		log.Fatalf("Failed to generate output: %v", err)
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
    subdomainx --wildcard <domains_file> [OPTIONS]

REQUIRED:
    --wildcard FILE        Path to file containing target domains (one per line)

OPTIONS:
    --version              Show version information
    --help                 Show this help message
    --check-tools          Check availability of enumeration tools
    --install-tools        Show installation instructions for missing tools
    
    # Output Options
    --name NAME            Unique name for output files (default: scan)
    --format FORMAT        Output format: json, txt, html (default: json)
    --output DIR           Output directory (default: output)
    
    # Performance Options
    --threads N            Number of concurrent threads (default: 10)
    --retries N            Number of retry attempts (default: 3)
    --timeout N            Timeout in seconds (default: 30)
    --rate-limit N         Rate limit per second (default: 100)
    
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
    
    # Configuration
    --config FILE          Use custom configuration file (optional)
    -v, --verbose          Enable verbose output

EXAMPLES:
    # Basic scan with all available tools
    subdomainx --wildcard domains.txt

    # Use only specific tools
    subdomainx --wildcard domains.txt --amass --subfinder --httpx

    # Generate HTML report
    subdomainx --wildcard domains.txt --format html --name my_scan

    # High-performance scan
    subdomainx --wildcard domains.txt --threads 20 --timeout 60

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
	// Validate wildcard file exists
	if !utils.FileExists(cfg.WildcardFile) {
		return fmt.Errorf("wildcard file not found: %s", cfg.WildcardFile)
	}

	// Validate output directory can be created
	if err := utils.EnsureDirectory(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Validate output format
	validFormats := map[string]bool{
		"json": true,
		"txt":  true,
		"html": true,
	}
	if !validFormats[cfg.OutputFormat] {
		return fmt.Errorf("invalid output format: %s. Supported formats: json, txt, html", cfg.OutputFormat)
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
