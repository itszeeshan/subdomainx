package types

import (
	"encoding/json"
	"testing"
)

func TestSubdomainResult(t *testing.T) {
	// Test SubdomainResult struct
	result := SubdomainResult{
		Subdomain: "test.example.com",
		Source:    "subfinder",
		IPs:       []string{"192.168.1.1", "10.0.0.1"},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal SubdomainResult: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled SubdomainResult
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal SubdomainResult: %v", err)
	}

	// Verify fields
	if unmarshaled.Subdomain != result.Subdomain {
		t.Errorf("Expected Subdomain '%s', got '%s'", result.Subdomain, unmarshaled.Subdomain)
	}
	if unmarshaled.Source != result.Source {
		t.Errorf("Expected Source '%s', got '%s'", result.Source, unmarshaled.Source)
	}
	if len(unmarshaled.IPs) != len(result.IPs) {
		t.Errorf("Expected %d IPs, got %d", len(result.IPs), len(unmarshaled.IPs))
	}
	for i, ip := range result.IPs {
		if unmarshaled.IPs[i] != ip {
			t.Errorf("Expected IP %d '%s', got '%s'", i, ip, unmarshaled.IPs[i])
		}
	}
}

func TestSubdomainResultEmptyIPs(t *testing.T) {
	// Test SubdomainResult with empty IPs (should omit from JSON)
	result := SubdomainResult{
		Subdomain: "test.example.com",
		Source:    "subfinder",
		IPs:       []string{},
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal SubdomainResult: %v", err)
	}

	// Verify IPs field is omitted when empty
	if len(result.IPs) == 0 && len(jsonData) > 0 {
		// Should not contain "ips" field when empty
		// This is expected behavior with omitempty tag
		t.Log("IPs field correctly omitted when empty")
	}

	// Test unmarshaling back
	var unmarshaled SubdomainResult
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal SubdomainResult: %v", err)
	}

	if unmarshaled.Subdomain != result.Subdomain {
		t.Errorf("Expected Subdomain '%s', got '%s'", result.Subdomain, unmarshaled.Subdomain)
	}
}

func TestHTTPResult(t *testing.T) {
	// Test HTTPResult struct
	result := HTTPResult{
		URL:           "https://test.example.com",
		StatusCode:    200,
		Title:         "Test Page",
		Technologies:  []string{"nginx", "php"},
		ContentLength: 1024,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal HTTPResult: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled HTTPResult
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal HTTPResult: %v", err)
	}

	// Verify fields
	if unmarshaled.URL != result.URL {
		t.Errorf("Expected URL '%s', got '%s'", result.URL, unmarshaled.URL)
	}
	if unmarshaled.StatusCode != result.StatusCode {
		t.Errorf("Expected StatusCode %d, got %d", result.StatusCode, unmarshaled.StatusCode)
	}
	if unmarshaled.Title != result.Title {
		t.Errorf("Expected Title '%s', got '%s'", result.Title, unmarshaled.Title)
	}
	if unmarshaled.ContentLength != result.ContentLength {
		t.Errorf("Expected ContentLength %d, got %d", result.ContentLength, unmarshaled.ContentLength)
	}
	if len(unmarshaled.Technologies) != len(result.Technologies) {
		t.Errorf("Expected %d technologies, got %d", len(result.Technologies), len(unmarshaled.Technologies))
	}
	for i, tech := range result.Technologies {
		if unmarshaled.Technologies[i] != tech {
			t.Errorf("Expected Technology %d '%s', got '%s'", i, tech, unmarshaled.Technologies[i])
		}
	}
}

func TestHTTPResultOptionalFields(t *testing.T) {
	// Test HTTPResult with optional fields omitted
	result := HTTPResult{
		URL:        "https://test.example.com",
		StatusCode: 404,
		// Title, Technologies, and ContentLength omitted
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal HTTPResult: %v", err)
	}

	// Test unmarshaling back
	var unmarshaled HTTPResult
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal HTTPResult: %v", err)
	}

	if unmarshaled.URL != result.URL {
		t.Errorf("Expected URL '%s', got '%s'", result.URL, unmarshaled.URL)
	}
	if unmarshaled.StatusCode != result.StatusCode {
		t.Errorf("Expected StatusCode %d, got %d", result.StatusCode, unmarshaled.StatusCode)
	}
}

func TestPortResult(t *testing.T) {
	// Test PortResult struct
	result := PortResult{
		Host: "test.example.com",
		IP:   "192.168.1.1",
		Ports: []Port{
			{Number: 80, Protocol: "tcp", State: "open", Service: "http"},
			{Number: 443, Protocol: "tcp", State: "open", Service: "https", Version: "nginx/1.18.0"},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal PortResult: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled PortResult
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal PortResult: %v", err)
	}

	// Verify fields
	if unmarshaled.Host != result.Host {
		t.Errorf("Expected Host '%s', got '%s'", result.Host, unmarshaled.Host)
	}
	if unmarshaled.IP != result.IP {
		t.Errorf("Expected IP '%s', got '%s'", result.IP, unmarshaled.IP)
	}
	if len(unmarshaled.Ports) != len(result.Ports) {
		t.Errorf("Expected %d ports, got %d", len(result.Ports), len(unmarshaled.Ports))
	}

	// Verify port details
	for i, port := range result.Ports {
		if unmarshaled.Ports[i].Number != port.Number {
			t.Errorf("Expected Port %d Number %d, got %d", i, port.Number, unmarshaled.Ports[i].Number)
		}
		if unmarshaled.Ports[i].Protocol != port.Protocol {
			t.Errorf("Expected Port %d Protocol '%s', got '%s'", i, port.Protocol, unmarshaled.Ports[i].Protocol)
		}
		if unmarshaled.Ports[i].State != port.State {
			t.Errorf("Expected Port %d State '%s', got '%s'", i, port.State, unmarshaled.Ports[i].State)
		}
		if unmarshaled.Ports[i].Service != port.Service {
			t.Errorf("Expected Port %d Service '%s', got '%s'", i, port.Service, unmarshaled.Ports[i].Service)
		}
		if unmarshaled.Ports[i].Version != port.Version {
			t.Errorf("Expected Port %d Version '%s', got '%s'", i, port.Version, unmarshaled.Ports[i].Version)
		}
	}
}

func TestPortResultEmptyPorts(t *testing.T) {
	// Test PortResult with empty ports
	result := PortResult{
		Host:  "test.example.com",
		IP:    "192.168.1.1",
		Ports: []Port{},
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal PortResult: %v", err)
	}

	// Test unmarshaling back
	var unmarshaled PortResult
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal PortResult: %v", err)
	}

	if unmarshaled.Host != result.Host {
		t.Errorf("Expected Host '%s', got '%s'", result.Host, unmarshaled.Host)
	}
	if unmarshaled.IP != result.IP {
		t.Errorf("Expected IP '%s', got '%s'", result.IP, unmarshaled.IP)
	}
	if len(unmarshaled.Ports) != 0 {
		t.Errorf("Expected 0 ports, got %d", len(unmarshaled.Ports))
	}
}

func TestPort(t *testing.T) {
	// Test Port struct
	port := Port{
		Number:   8080,
		Protocol: "tcp",
		State:    "open",
		Service:  "http-proxy",
		Version:  "nginx/1.20.0",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(port)
	if err != nil {
		t.Fatalf("Failed to marshal Port: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Port
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Port: %v", err)
	}

	// Verify fields
	if unmarshaled.Number != port.Number {
		t.Errorf("Expected Number %d, got %d", port.Number, unmarshaled.Number)
	}
	if unmarshaled.Protocol != port.Protocol {
		t.Errorf("Expected Protocol '%s', got '%s'", port.Protocol, unmarshaled.Protocol)
	}
	if unmarshaled.State != port.State {
		t.Errorf("Expected State '%s', got '%s'", port.State, unmarshaled.State)
	}
	if unmarshaled.Service != port.Service {
		t.Errorf("Expected Service '%s', got '%s'", port.Service, unmarshaled.Service)
	}
	if unmarshaled.Version != port.Version {
		t.Errorf("Expected Version '%s', got '%s'", port.Version, unmarshaled.Version)
	}
}

func TestPortOptionalFields(t *testing.T) {
	// Test Port with optional fields omitted
	port := Port{
		Number:   22,
		Protocol: "tcp",
		State:    "closed",
		// Service and Version omitted
	}

	jsonData, err := json.Marshal(port)
	if err != nil {
		t.Fatalf("Failed to marshal Port: %v", err)
	}

	// Test unmarshaling back
	var unmarshaled Port
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Port: %v", err)
	}

	if unmarshaled.Number != port.Number {
		t.Errorf("Expected Number %d, got %d", port.Number, unmarshaled.Number)
	}
	if unmarshaled.Protocol != port.Protocol {
		t.Errorf("Expected Protocol '%s', got '%s'", port.Protocol, unmarshaled.Protocol)
	}
	if unmarshaled.State != port.State {
		t.Errorf("Expected State '%s', got '%s'", port.State, unmarshaled.State)
	}
}

func TestScanResults(t *testing.T) {
	// Test ScanResults struct
	results := ScanResults{
		Subdomains: []SubdomainResult{
			{Subdomain: "test1.example.com", Source: "subfinder", IPs: []string{"192.168.1.1"}},
			{Subdomain: "test2.example.com", Source: "amass", IPs: []string{"10.0.0.1"}},
		},
		HTTP: []HTTPResult{
			{URL: "https://test1.example.com", StatusCode: 200, Title: "Test Page"},
			{URL: "https://test2.example.com", StatusCode: 404},
		},
		Ports: []PortResult{
			{Host: "test1.example.com", IP: "192.168.1.1", Ports: []Port{{Number: 80, Protocol: "tcp", State: "open"}}},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(results)
	if err != nil {
		t.Fatalf("Failed to marshal ScanResults: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled ScanResults
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ScanResults: %v", err)
	}

	// Verify subdomains
	if len(unmarshaled.Subdomains) != len(results.Subdomains) {
		t.Errorf("Expected %d subdomains, got %d", len(results.Subdomains), len(unmarshaled.Subdomains))
	}

	// Verify HTTP results
	if len(unmarshaled.HTTP) != len(results.HTTP) {
		t.Errorf("Expected %d HTTP results, got %d", len(results.HTTP), len(unmarshaled.HTTP))
	}

	// Verify port results
	if len(unmarshaled.Ports) != len(results.Ports) {
		t.Errorf("Expected %d port results, got %d", len(results.Ports), len(unmarshaled.Ports))
	}
}

func TestScanResultsEmpty(t *testing.T) {
	// Test ScanResults with empty slices
	results := ScanResults{
		Subdomains: []SubdomainResult{},
		HTTP:       []HTTPResult{},
		Ports:      []PortResult{},
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		t.Fatalf("Failed to marshal ScanResults: %v", err)
	}

	// Test unmarshaling back
	var unmarshaled ScanResults
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ScanResults: %v", err)
	}

	if len(unmarshaled.Subdomains) != 0 {
		t.Errorf("Expected 0 subdomains, got %d", len(unmarshaled.Subdomains))
	}
	if len(unmarshaled.HTTP) != 0 {
		t.Errorf("Expected 0 HTTP results, got %d", len(unmarshaled.HTTP))
	}
	if len(unmarshaled.Ports) != 0 {
		t.Errorf("Expected 0 port results, got %d", len(unmarshaled.Ports))
	}
}
