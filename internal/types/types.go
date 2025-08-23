package types

type SubdomainResult struct {
	Subdomain string   `json:"subdomain"`
	Source    string   `json:"source"`
	IPs       []string `json:"ips,omitempty"`
}

type HTTPResult struct {
	URL           string   `json:"url"`
	StatusCode    int      `json:"status_code"`
	Title         string   `json:"title,omitempty"`
	Technologies  []string `json:"technologies,omitempty"`
	ContentLength int      `json:"content_length,omitempty"`
}

type PortResult struct {
	Host  string `json:"host"`
	IP    string `json:"ip,omitempty"`
	Ports []Port `json:"ports,omitempty"`
}

type Port struct {
	Number   int    `json:"number"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	Service  string `json:"service,omitempty"`
	Version  string `json:"version,omitempty"`
}

type ScanResults struct {
	Subdomains []SubdomainResult `json:"subdomains"`
	HTTP       []HTTPResult      `json:"http,omitempty"`
	Ports      []PortResult      `json:"ports,omitempty"`
}
