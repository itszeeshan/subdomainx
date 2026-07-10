package notify

import (
	"fmt"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/v2/internal/diff"
)

// ScanSummary is the payload delivered to every notifier.
type ScanSummary struct {
	ScanID          string
	Domain          string
	TotalSubdomains int
	TotalHTTP       int
	TotalPorts      int
	Duration        time.Duration
	Error           string           // non-empty if scan failed
	Diff            *diff.DiffResult // nil when --diff not used
}

// Notifier sends a scan summary to one notification channel.
type Notifier interface {
	Name() string
	Send(summary ScanSummary) error
}

var notifiers = map[string]Notifier{}

// Register adds a notifier to the registry. Called by each channel's init().
func Register(n Notifier) {
	notifiers[n.Name()] = n
}

// ValidChannels returns the list of supported channel names.
func ValidChannels() []string {
	channels := make([]string, 0, len(notifiers))
	for name := range notifiers {
		channels = append(channels, name)
	}
	return channels
}

// IsValidChannel checks if a channel name is registered.
func IsValidChannel(name string) bool {
	_, ok := notifiers[name]
	return ok
}

// Send dispatches the summary to all requested channels. Errors from
// individual channels are collected and returned as a combined error.
// Partial success is expected — one channel failing does not block others.
func Send(channels []string, summary ScanSummary) error {
	var errs []string
	for _, ch := range channels {
		n, ok := notifiers[ch]
		if !ok {
			errs = append(errs, fmt.Sprintf("%s: unknown channel", ch))
			continue
		}
		if err := n.Send(summary); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", ch, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("notification errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
