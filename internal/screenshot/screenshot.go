package screenshot

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/itszeeshan/subdomainx/v2/internal/config"
	"github.com/itszeeshan/subdomainx/v2/internal/types"
	"github.com/itszeeshan/subdomainx/v2/internal/utils"
)

// CaptureAll takes screenshots of all HTTP-alive subdomains concurrently.
// Returns the number of successfully captured screenshots.
func CaptureAll(cfg *config.Config, httpResults []types.HTTPResult) (int, error) {
	width, height := parseResolution(cfg.ScreenshotResolution)
	timeout := cfg.ScreenshotTimeout
	if timeout <= 0 {
		timeout = 10
	}

	dir := cfg.ScreenshotDir
	if dir == "" {
		dir = filepath.Join(cfg.OutputDir, "screenshots")
		cfg.ScreenshotDir = dir // persist so HTML report can find screenshots
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create screenshot directory: %w", err)
	}

	// Deduplicate: prefer HTTPS over HTTP for the same host.
	targets := deduplicateTargets(httpResults)

	// Limit to MaxHTTPTargets.
	if cfg.MaxHTTPTargets > 0 && len(targets) > cfg.MaxHTTPTargets {
		targets = targets[:cfg.MaxHTTPTargets]
	}

	if len(targets) == 0 {
		return 0, nil
	}

	// Create a headless Chrome allocator context shared across all tabs.
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.WindowSize(width, height),
			chromedp.Flag("ignore-certificate-errors", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
		)...,
	)
	defer allocCancel()

	// Create a browser context from the allocator.
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)
	defer browserCancel()

	// Start the browser by running an empty task.
	if err := chromedp.Run(browserCtx); err != nil {
		return 0, fmt.Errorf("failed to start browser: %w", err)
	}

	pool := utils.NewWorkerPool(cfg.Threads, cfg.RateLimit)
	defer pool.Stop()

	var captured atomic.Int32
	var wg sync.WaitGroup

	for _, target := range targets {
		target := target
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()

			outPath := filepath.Join(dir, sanitizeFilename(target)+".png")
			if err := captureSingle(browserCtx, target, outPath, width, height, timeout); err != nil {
				log.Printf("Screenshot failed for %s: %v", target, err)
				return
			}
			captured.Add(1)
		})
	}

	wg.Wait()
	return int(captured.Load()), nil
}

func captureSingle(browserCtx context.Context, targetURL, outputPath string, width, height, timeoutSec int) error {
	ctx, cancel := chromedp.NewContext(browserCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(int64(width), int64(height)),
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body"),
		chromedp.CaptureScreenshot(&buf),
	)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, buf, 0644)
}

// deduplicateTargets picks HTTPS over HTTP when both exist for the same host.
func deduplicateTargets(results []types.HTTPResult) []string {
	seen := make(map[string]string) // host -> best URL
	for _, r := range results {
		parsed, err := url.Parse(r.URL)
		if err != nil {
			continue
		}
		host := parsed.Hostname()
		existing, exists := seen[host]
		if !exists {
			seen[host] = r.URL
		} else if parsed.Scheme == "https" && strings.HasPrefix(existing, "http://") {
			seen[host] = r.URL
		}
	}

	targets := make([]string, 0, len(seen))
	for _, u := range seen {
		targets = append(targets, u)
	}
	return targets
}

func sanitizeFilename(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return strings.NewReplacer("://", "_", "/", "_", ":", "_").Replace(rawURL)
	}
	host := parsed.Hostname()
	if port := parsed.Port(); port != "" && port != "80" && port != "443" {
		host = host + "_" + port
	}
	return host
}

func parseResolution(res string) (int, int) {
	parts := strings.SplitN(strings.ToLower(res), "x", 2)
	if len(parts) == 2 {
		w, err1 := strconv.Atoi(parts[0])
		h, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil && w > 0 && h > 0 {
			return w, h
		}
	}
	return 1280, 720
}
