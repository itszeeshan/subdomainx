package notify

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func init() {
	Register(&emailNotifier{})
}

type emailNotifier struct{}

func (e *emailNotifier) Name() string { return "email" }

func (e *emailNotifier) Send(summary ScanSummary) error {
	host := os.Getenv("SUBDOMAINX_SMTP_HOST")
	port := os.Getenv("SUBDOMAINX_SMTP_PORT")
	user := os.Getenv("SUBDOMAINX_SMTP_USER")
	pass := os.Getenv("SUBDOMAINX_SMTP_PASS")
	to := os.Getenv("SUBDOMAINX_NOTIFY_EMAIL")

	if host == "" {
		return fmt.Errorf("SUBDOMAINX_SMTP_HOST environment variable not set")
	}
	if to == "" {
		return fmt.Errorf("SUBDOMAINX_NOTIFY_EMAIL environment variable not set")
	}
	if port == "" {
		port = "587"
	}

	subject := fmt.Sprintf("SubdomainX: Scan completed for %s", summary.Domain)
	if summary.Error != "" {
		subject = fmt.Sprintf("SubdomainX: Scan FAILED for %s", summary.Domain)
	}

	recipients := strings.Split(to, ",")
	body := FormatPlainText(summary)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		user, to, subject, body)

	addr := host + ":" + port
	var auth smtp.Auth
	if user != "" && pass != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}

	if err := smtp.SendMail(addr, auth, user, recipients, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
