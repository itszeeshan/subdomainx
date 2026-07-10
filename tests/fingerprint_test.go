package tests

import (
	"net/http"
	"testing"

	"github.com/itszeeshan/subdomainx/v2/internal/scanner"
)

func TestFingerprintTechnologies_Headers(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Server":       {"nginx/1.24.0"},
			"X-Powered-By": {"Express"},
		},
	}

	techs := scanner.FingerprintTechnologies(resp, nil)

	found := make(map[string]bool)
	for _, tech := range techs {
		found[tech.Name] = true
	}

	if !found["nginx"] {
		t.Error("Expected to detect nginx from Server header")
	}
	if !found["Express"] {
		t.Error("Expected to detect Express from X-Powered-By header")
	}

	// Verify nginx has version
	for _, tech := range techs {
		if tech.Name == "nginx" {
			if tech.Version != "1.24.0" {
				t.Errorf("Expected nginx version 1.24.0, got %q", tech.Version)
			}
			if tech.Category != "Web Server" {
				t.Errorf("Expected category Web Server, got %q", tech.Category)
			}
		}
	}
}

func TestFingerprintTechnologies_Cookies(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Set-Cookie": {
				"PHPSESSID=abc123; path=/",
				"csrftoken=xyz789; path=/",
			},
		},
	}

	techs := scanner.FingerprintTechnologies(resp, nil)

	found := make(map[string]string)
	for _, tech := range techs {
		found[tech.Name] = tech.Category
	}

	if _, ok := found["PHP"]; !ok {
		t.Error("Expected to detect PHP from PHPSESSID cookie")
	}
	if _, ok := found["Django"]; !ok {
		t.Error("Expected to detect Django from csrftoken cookie")
	}
}

func TestFingerprintTechnologies_CDN(t *testing.T) {
	header := make(http.Header)
	header.Set("CF-RAY", "abc123-SJC")
	resp := &http.Response{
		Header: header,
	}

	techs := scanner.FingerprintTechnologies(resp, nil)

	found := false
	for _, tech := range techs {
		if tech.Name == "Cloudflare" && tech.Category == "CDN" {
			found = true
		}
	}
	if !found {
		t.Error("Expected to detect Cloudflare from CF-RAY header")
	}
}

func TestFingerprintTechnologies_HTMLBody(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}

	body := []byte(`<!DOCTYPE html>
<html>
<head>
	<meta name="generator" content="WordPress 6.4.2">
	<link rel="stylesheet" href="/wp-content/themes/default/style.css">
	<script src="https://cdn.example.com/jquery-3.7.1.min.js"></script>
	<script src="https://www.googletagmanager.com/gtag/js?id=G-TEST"></script>
</head>
<body>
	<div id="__next">Next.js app</div>
</body>
</html>`)

	techs := scanner.FingerprintTechnologies(resp, body)

	found := make(map[string]string)
	for _, tech := range techs {
		found[tech.Name] = tech.Version
	}

	if _, ok := found["WordPress"]; !ok {
		t.Error("Expected to detect WordPress from meta generator and wp-content")
	}
	if v, ok := found["jQuery"]; !ok {
		t.Error("Expected to detect jQuery from script src")
	} else if v != "3.7.1" {
		t.Errorf("Expected jQuery version 3.7.1, got %q", v)
	}
	if _, ok := found["Google Analytics"]; !ok {
		t.Error("Expected to detect Google Analytics from googletagmanager.com")
	}
	if _, ok := found["Next.js"]; !ok {
		t.Error("Expected to detect Next.js from __next div")
	}
}

func TestFingerprintTechnologies_NoDuplicates(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Server": {"nginx/1.24.0"},
		},
	}

	// Body also references nginx
	body := []byte(`<meta name="generator" content="WordPress 6.4">`)

	techs := scanner.FingerprintTechnologies(resp, body)

	counts := make(map[string]int)
	for _, tech := range techs {
		counts[tech.Name]++
	}

	for name, count := range counts {
		if count > 1 {
			t.Errorf("Technology %q detected %d times, expected 1", name, count)
		}
	}
}

func TestFingerprintTechnologies_EmptyResponse(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}

	techs := scanner.FingerprintTechnologies(resp, nil)

	if len(techs) != 0 {
		t.Errorf("Expected 0 technologies for empty response, got %d", len(techs))
	}
}
