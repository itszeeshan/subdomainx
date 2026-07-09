package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func init() {
	Register(&telegramNotifier{})
}

type telegramNotifier struct{}

func (t *telegramNotifier) Name() string { return "telegram" }

func (t *telegramNotifier) Send(summary ScanSummary) error {
	token := os.Getenv("SUBDOMAINX_TELEGRAM_TOKEN")
	chatID := os.Getenv("SUBDOMAINX_TELEGRAM_CHAT_ID")
	if token == "" {
		return fmt.Errorf("SUBDOMAINX_TELEGRAM_TOKEN environment variable not set")
	}
	if chatID == "" {
		return fmt.Errorf("SUBDOMAINX_TELEGRAM_CHAT_ID environment variable not set")
	}

	payload := map[string]string{
		"chat_id": chatID,
		"text":    FormatPlainText(summary),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal telegram payload: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send telegram notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}
	return nil
}
