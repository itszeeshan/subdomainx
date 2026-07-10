package tests

import (
	"testing"

	"github.com/itszeeshan/subdomainx/v2/internal/scanner"
)

func TestMatchesCnamePattern(t *testing.T) {
	tests := []struct {
		cname    string
		patterns []string
		want     bool
	}{
		{"old-app.herokuapp.com", []string{".herokuapp.com"}, true},
		{"example.github.io", []string{".github.io"}, true},
		{"mybucket.s3.amazonaws.com", []string{".s3.amazonaws.com", ".s3-website"}, true},
		{"example.com", []string{".herokuapp.com"}, false},
		{"herokuapp.com.evil.com", []string{".herokuapp.com"}, false},
		{"test.netlify.app", []string{".netlify.app", ".netlify.com"}, true},
		{"test.netlify.com", []string{".netlify.app", ".netlify.com"}, true},
		{"test.azurewebsites.net", []string{".azurewebsites.net", ".cloudapp.net"}, true},
	}

	for _, tt := range tests {
		got := scanner.MatchesCnamePattern(tt.cname, tt.patterns)
		if got != tt.want {
			t.Errorf("MatchesCnamePattern(%q, %v) = %v, want %v", tt.cname, tt.patterns, got, tt.want)
		}
	}
}

func TestExtractHostFromURL(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://example.com", "example.com"},
		{"http://example.com:8080/path", "example.com"},
		{"https://sub.example.com/", "sub.example.com"},
		{"http://test.example.com:443/path?query=1", "test.example.com"},
	}

	for _, tt := range tests {
		got := scanner.ExtractHostFromURL(tt.url)
		if got != tt.want {
			t.Errorf("ExtractHostFromURL(%q) = %q, want %q", tt.url, got, tt.want)
		}
	}
}
