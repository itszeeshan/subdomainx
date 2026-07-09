package types

type SubdomainResult struct {
	Subdomain string   `json:"subdomain"`
	Source    string   `json:"source"`
	IPs       []string `json:"ips,omitempty"`
}

type LinkHeader struct {
	URL        string   `json:"url"`
	Rel        string   `json:"rel"`
	Subdomains []string `json:"subdomains,omitempty"`
}

type Technology struct {
	Name     string `json:"name"`
	Version  string `json:"version,omitempty"`
	Category string `json:"category"`
}

type HTTPResult struct {
	URL           string       `json:"url"`
	StatusCode    int          `json:"status_code"`
	Title         string       `json:"title,omitempty"`
	Technologies  []string     `json:"technologies,omitempty"`
	DetectedTech  []Technology `json:"detected_tech,omitempty"`
	ContentLength int          `json:"content_length,omitempty"`
	LinkHeaders   []LinkHeader `json:"link_headers,omitempty"`
}

type TakeoverResult struct {
	Subdomain string `json:"subdomain"`
	CNAME     string `json:"cname,omitempty"`
	Risk      string `json:"risk"`
	Service   string `json:"service"`
	Evidence  string `json:"evidence"`
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

// WaybackEntry holds historical URLs discovered for a subdomain.
type WaybackEntry struct {
	Subdomain string   `json:"subdomain"`
	Domain    string   `json:"domain"`
	URLs      []string `json:"urls"`
}

type ScanResults struct {
	Subdomains []SubdomainResult `json:"subdomains"`
	HTTP       []HTTPResult      `json:"http,omitempty"`
	Ports      []PortResult      `json:"ports,omitempty"`
	Wayback    []WaybackEntry    `json:"wayback,omitempty"`
	Takeover   []TakeoverResult  `json:"takeover,omitempty"`
}
