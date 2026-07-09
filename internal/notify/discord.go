package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const discordMaxLength = 2000

func init() {
	Register(&discordNotifier{})
}

type discordNotifier struct{}

func (d *discordNotifier) Name() string { return "discord" }

func (d *discordNotifier) Send(summary ScanSummary) error {
	webhook := os.Getenv("SUBDOMAINX_DISCORD_WEBHOOK")
	if webhook == "" {
		return fmt.Errorf("SUBDOMAINX_DISCORD_WEBHOOK environment variable not set")
	}

	content := FormatMarkdown(summary)
	if len(content) > discordMaxLength {
		content = content[:discordMaxLength-20] + "\n...(truncated)"
	}

	payload := map[string]string{
		"content": content,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal discord payload: %w", err)
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send discord notification: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}
	return nil
}
