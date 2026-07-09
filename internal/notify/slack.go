package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func init() {
	Register(&slackNotifier{})
}

type slackNotifier struct{}

func (s *slackNotifier) Name() string { return "slack" }

func (s *slackNotifier) Send(summary ScanSummary) error {
	webhook := os.Getenv("SUBDOMAINX_SLACK_WEBHOOK")
	if webhook == "" {
		return fmt.Errorf("SUBDOMAINX_SLACK_WEBHOOK environment variable not set")
	}

	payload := map[string]string{
		"text": FormatMarkdown(summary),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack webhook returned status %d", resp.StatusCode)
	}
	return nil
}
