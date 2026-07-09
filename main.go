package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

func main() {
	// ---- Flag definitions ----
	var (
		showVersion     = flag.Bool("version", false, "Show version information")
		showHelp        = flag.Bool("help", false, "Show help information")
		checkTools      = flag.Bool("check-tools", false, "Check tool availability")
		installTools    = flag.Bool("install-tools", false, "Show tool installation instructions")
		configFile      = flag.String("config", "", "Path to configuration file (optional)")
		wildcardFile    = flag.String("wildcard", "", "Path to wildcard file containing domains")
		uniqueName      = flag.String("name", "scan", "Unique name for output files")
		outputFormat    = flag.String("format", "", "Output format (json, txt, html, csv, burp, nessus, zap)")
		outputDir       = flag.String("output", "output", "Output directory")
		threads         = flag.Int("threads", 10, "Number of threads")
		retries         = flag.Int("retries", 3, "Number of retry attempts")
		timeout         = flag.Int("timeout", 30, "Timeout in seconds")
		rateLimit       = flag.Int("rate-limit", 100, "Rate limit per second")
		wordlist        = flag.String("wordlist", "", "Custom wordlist file for brute-forcing")
		resume          = flag.String("resume", "", "Resume scan from checkpoint (scan ID)")
		listCheckpoints = flag.Bool("list-checkpoints", false, "List available checkpoints")
		verbose         = flag.Bool("verbose", false, "Verbose output")
		statusCodes     = flag.String("status-codes", "", "Filter by HTTP status codes (e.g., '200,301,302')")
		ports           = flag.String("ports", "", "Filter by ports (e.g., '80,443,8080')")
		maxHTTPTargets  = flag.Int("max-http-targets", 1000, "Maximum number of subdomains to scan with httpx")
		diffMode        = flag.Bool("diff", false, "Compare results against previous scan")
		baselineFile    = flag.String("baseline", "", "Baseline results file for diff comparison")
		notifyFlag      = flag.String("notify", "", "Notification channels (comma-separated: slack,discord,telegram,email)")

		flags = toolFlags{}
	)

	flag.BoolVar(&flags.useSubfinder, "subfinder", false, "Use subfinder tool")
	flag.BoolVar(&flags.useAmass, "amass", false, "Use amass tool")
	flag.BoolVar(&flags.useFindomain, "findomain", false, "Use findomain tool")
	flag.BoolVar(&flags.useAssetfinder, "assetfinder", false, "Use assetfinder tool")
	flag.BoolVar(&flags.useSublist3r, "sublist3r", false, "Use sublist3r tool")
	flag.BoolVar(&flags.useKnockpy, "knockpy", false, "Use knockpy tool")
	flag.BoolVar(&flags.useDnsrecon, "dnsrecon", false, "Use dnsrecon tool")
	flag.BoolVar(&flags.useFierce, "fierce", false, "Use fierce tool")
	flag.BoolVar(&flags.useMassdns, "massdns", false, "Use massdns tool")
	flag.BoolVar(&flags.useAltdns, "altdns", false, "Use altdns tool")
	flag.BoolVar(&flags.useSecurityTrails, "securitytrails", false, "Use SecurityTrails API")
	flag.BoolVar(&flags.useVirusTotal, "virustotal", false, "Use VirusTotal API")
	flag.BoolVar(&flags.useCensys, "censys", false, "Use Censys API")
	flag.BoolVar(&flags.useCrtSh, "crtsh", false, "Use crt.sh Certificate Transparency API")
	flag.BoolVar(&flags.useURLScan, "urlscan", false, "Use URLScan.io API")
	flag.BoolVar(&flags.useHackerTarget, "hackertarget", false, "Use HackerTarget API")
	flag.BoolVar(&flags.useWaybackURLs, "waybackurls", false, "Use waybackurls tool")
	flag.BoolVar(&flags.useLinkHeader, "linkheader", false, "Use Link Header enumeration")
	flag.BoolVar(&flags.useHttpx, "httpx", false, "Use httpx for HTTP scanning")
	flag.BoolVar(&flags.useSmap, "smap", false, "Use smap for port scanning")

	flag.Parse()

	// ---- Early-exit commands ----
	if *showVersion {
		fmt.Println("SubdomainX v1.5.0")
		return
	}
	if *showHelp {
		showUsage()
		return
	}
	if *checkTools {
		utils.DisplayToolStatus()
		return
	}
	if *installTools {
		if err := utils.InstallTools(); err != nil {
			log.Fatalf("Installation failed: %v", err)
		}
		return
	}
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

	showBanner()

	if *verbose {
		utils.StartResourceMonitoring()
		defer utils.StopResourceMonitoring()
	}

	available, missing := utils.CheckAllTools()
	if len(missing) > 0 && *verbose {
		fmt.Printf("⚠️  Warning: %d tools are missing. SubdomainX will work with available tools.\n", len(missing))
		fmt.Println("💡 Use --install-tools to see installation instructions.")
		fmt.Printf("✅ Available tools: %d\n\n", len(available))
	}

	// ---- Build base config from CLI flags ----
	args := flag.Args()
	hasDomainArg := len(args) > 0

	cfg := &config.Config{
		OutputDir:      *outputDir,
		OutputFormat:   "json",
		Threads:        *threads,
		Retries:        *retries,
		Timeout:        *timeout,
		RateLimit:      *rateLimit,
		MaxHTTPTargets: *maxHTTPTargets,
		Tools:          make(map[string]bool),
		Filters:        make(map[string]string),
	}
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
	if *statusCodes != "" {
		cfg.Filters["status_code"] = *statusCodes
	}
	if *ports != "" {
		cfg.Filters["ports"] = *ports
	}
	cfg.DiffEnabled = *diffMode
	if *baselineFile != "" {
		cfg.BaselineFile = *baselineFile
		cfg.DiffEnabled = true // --baseline implies --diff
	}
	if *notifyFlag != "" {
		cfg.NotifyChannels = strings.Split(*notifyFlag, ",")
	}

	if cfg.WildcardFile == "" && len(args) == 0 && *resume == "" {
		log.Fatalf("Error: Either --wildcard file, a domain argument, or --resume is required. Use --help for usage information.")
	}

	// ---- Domain input: create temp file when a bare domain is given ----
	cleanup, err := setupWildcardFile(cfg, args)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer cleanup()

	// ---- Config file merging (CLI wins) ----
	cfg = loadAndMergeConfig(cfg, *configFile, hasDomainArg, *resume)

	// ---- Tool selection (after config merge so CLI always wins) ----
	applyToolSelection(cfg, flags, *verbose)

	// ---- Validate and create output directory ----
	if err := validateCLIInput(cfg); err != nil {
		log.Fatalf("Validation error: %v", err)
	}
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// ---- Checkpoint init and scan pipeline ----
	state, err := initScanState(cfg, args, *resume, *outputDir)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err := executeScanPipeline(cfg, state, *resume); err != nil {
		log.Fatalf("%v", err)
	}
}
