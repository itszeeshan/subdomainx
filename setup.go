package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

// toolFlags holds all per-tool CLI flags so they can be passed as one argument.
type toolFlags struct {
	useSubfinder      bool
	useAmass          bool
	useFindomain      bool
	useAssetfinder    bool
	useSublist3r      bool
	useKnockpy        bool
	useDnsrecon       bool
	useFierce         bool
	useMassdns        bool
	useAltdns         bool
	useSecurityTrails bool
	useVirusTotal     bool
	useCensys         bool
	useCrtSh          bool
	useURLScan        bool
	useHackerTarget   bool
	useWaybackURLs    bool
	useLinkHeader     bool
	useHttpx          bool
	useSmap           bool
}

func (f toolFlags) anySelected() bool {
	return f.useSubfinder || f.useAmass || f.useFindomain || f.useAssetfinder ||
		f.useSublist3r || f.useKnockpy || f.useDnsrecon || f.useFierce ||
		f.useMassdns || f.useAltdns || f.useSecurityTrails || f.useVirusTotal ||
		f.useCensys || f.useCrtSh || f.useURLScan ||
		f.useHackerTarget || f.useWaybackURLs || f.useLinkHeader ||
		f.useHttpx || f.useSmap
}

// setupWildcardFile handles the single-domain argument: validates the domain,
// writes it to a temp file, and sets cfg.WildcardFile. The returned cleanup
// function removes the temp file and should be deferred by the caller.
func setupWildcardFile(cfg *config.Config, args []string) (cleanup func(), err error) {
	cleanup = func() {} // no-op by default
	if cfg.WildcardFile != "" || len(args) == 0 {
		return cleanup, nil
	}

	domain := args[0]
	if err := utils.ValidateDomain(domain); err != nil {
		return cleanup, fmt.Errorf("invalid domain '%s': %v", domain, err)
	}

	tmpFile, err := os.CreateTemp("", "subdomainx_domain_*.txt")
	if err != nil {
		return cleanup, fmt.Errorf("failed to create temporary file: %v", err)
	}

	if _, err := tmpFile.WriteString(domain + "\n"); err != nil {
		os.Remove(tmpFile.Name())
		return cleanup, fmt.Errorf("failed to write domain to temporary file: %v", err)
	}
	tmpFile.Close()

	cfg.WildcardFile = tmpFile.Name()
	if cfg.UniqueName == "scan" {
		cfg.UniqueName = domain
	}

	return func() { os.Remove(tmpFile.Name()) }, nil
}

// loadAndMergeConfig loads a config file (or the default) and merges it with
// the CLI-built cfg, returning the merged result. CLI values win.
func loadAndMergeConfig(cfg *config.Config, configFile string, hasDomainArg bool, resume string) *config.Config {
	if configFile != "" {
		fileCfg, err := config.LoadConfigFromFile(configFile)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
		return mergeConfig(fileCfg, cfg)
	}
	if !hasDomainArg && resume == "" {
		if defaultCfg, err := config.LoadConfig(); err == nil {
			return mergeConfig(defaultCfg, cfg)
		}
	}
	return cfg
}

// applyToolSelection writes the enabled/disabled tool map into cfg based on
// the CLI flags. When specific tools are chosen those are applied exclusively;
// otherwise every tool is enabled. Must be called after config file merging.
func applyToolSelection(cfg *config.Config, flags toolFlags, verbose bool) {
	if flags.anySelected() {
		cfg.Tools = map[string]bool{
			"subfinder":      flags.useSubfinder,
			"amass":          flags.useAmass,
			"findomain":      flags.useFindomain,
			"assetfinder":    flags.useAssetfinder,
			"sublist3r":      flags.useSublist3r,
			"knockpy":        flags.useKnockpy,
			"dnsrecon":       flags.useDnsrecon,
			"fierce":         flags.useFierce,
			"massdns":        flags.useMassdns,
			"altdns":         flags.useAltdns,
			"securitytrails": flags.useSecurityTrails,
			"virustotal":     flags.useVirusTotal,
			"censys":         flags.useCensys,
			"crtsh":          flags.useCrtSh,
			"urlscan":        flags.useURLScan,
			"hackertarget":   flags.useHackerTarget,
			"waybackurls":    flags.useWaybackURLs,
			"linkheader":     flags.useLinkHeader,
			"httpx":          flags.useHttpx,
			"smap":           flags.useSmap,
		}
		if verbose {
			var selected []string
			for tool, enabled := range cfg.Tools {
				if enabled {
					selected = append(selected, tool)
				}
			}
			fmt.Printf("🔧 CLI Tools selected: %s\n", strings.Join(selected, ", "))
		}
		return
	}

	// No specific tools chosen — enable everything.
	if cfg.Tools == nil {
		cfg.Tools = make(map[string]bool)
	}
	for _, tool := range []string{
		"subfinder", "amass", "findomain", "assetfinder", "sublist3r",
		"knockpy", "dnsrecon", "fierce", "massdns", "altdns",
		"securitytrails", "virustotal", "censys", "crtsh", "urlscan",
		"hackertarget", "waybackurls", "linkheader",
		"httpx", "smap",
	} {
		cfg.Tools[tool] = true
	}
}

// mergeConfig merges two configs with cfg2 taking precedence over cfg1.
func mergeConfig(cfg1, cfg2 *config.Config) *config.Config {
	result := &config.Config{
		WildcardFile:   cfg1.WildcardFile,
		UniqueName:     cfg1.UniqueName,
		OutputDir:      cfg1.OutputDir,
		OutputFormat:   cfg1.OutputFormat,
		Threads:        cfg1.Threads,
		Retries:        cfg1.Retries,
		Timeout:        cfg1.Timeout,
		RateLimit:      cfg1.RateLimit,
		Wordlist:       cfg1.Wordlist,
		MaxHTTPTargets: cfg1.MaxHTTPTargets,
		Tools:          make(map[string]bool),
		Filters:        make(map[string]string),
	}

	for k, v := range cfg1.Tools {
		result.Tools[k] = v
	}
	for k, v := range cfg1.Filters {
		result.Filters[k] = v
	}

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
	if cfg2.MaxHTTPTargets > 0 {
		result.MaxHTTPTargets = cfg2.MaxHTTPTargets
	}

	if cfg2.DiffEnabled {
		result.DiffEnabled = true
	}
	if cfg2.BaselineFile != "" {
		result.BaselineFile = cfg2.BaselineFile
	}
	if len(cfg2.NotifyChannels) > 0 {
		result.NotifyChannels = cfg2.NotifyChannels
	}

	for k, v := range cfg2.Tools {
		result.Tools[k] = v
	}
	for k, v := range cfg2.Filters {
		result.Filters[k] = v
	}

	return result
}

// validateCLIInput validates the resolved configuration before scanning begins.
func validateCLIInput(cfg *config.Config) error {
	// Skip file-existence check for the temp file created for single-domain mode.
	if cfg.WildcardFile != "" &&
		!strings.Contains(cfg.WildcardFile, "subdomainx_domain_") &&
		!utils.FileExists(cfg.WildcardFile) {
		return fmt.Errorf("wildcard file not found: %s", cfg.WildcardFile)
	}

	if err := utils.EnsureDirectory(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	validFormats := map[string]bool{
		"json": true, "txt": true, "html": true,
		"zap": true, "burp": true, "nessus": true, "csv": true,
	}
	if !validFormats[cfg.OutputFormat] {
		return fmt.Errorf("invalid output format: %s. Supported: json, txt, html, zap, burp, nessus, csv", cfg.OutputFormat)
	}

	if cfg.Threads <= 0 {
		return fmt.Errorf("threads must be greater than 0")
	}
	if cfg.Retries < 0 {
		return fmt.Errorf("retries cannot be negative")
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}
	if cfg.RateLimit <= 0 {
		return fmt.Errorf("rate limit must be greater than 0")
	}
	if cfg.Wordlist != "" && !utils.FileExists(cfg.Wordlist) {
		return fmt.Errorf("wordlist file not found: %s", cfg.Wordlist)
	}
	if cfg.BaselineFile != "" && !utils.FileExists(cfg.BaselineFile) {
		return fmt.Errorf("baseline file not found: %s", cfg.BaselineFile)
	}
	validChannels := map[string]bool{"slack": true, "discord": true, "telegram": true, "email": true}
	for _, ch := range cfg.NotifyChannels {
		if !validChannels[ch] {
			return fmt.Errorf("invalid notification channel: %s. Supported: slack, discord, telegram, email", ch)
		}
	}

	return nil
}
