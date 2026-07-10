package scanner

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/itszeeshan/subdomainx/v2/internal/types"
)

// headerFingerprint maps an HTTP header to detection rules.
type headerFingerprint struct {
	Header   string
	Category string
	// If Name is set, use it as-is. Otherwise derive name from header value.
	Name string
}

var headerFingerprints = []headerFingerprint{
	{Header: "Server", Category: "Web Server"},
	{Header: "X-Powered-By", Category: "Framework"},
	{Header: "X-AspNet-Version", Category: "Framework", Name: "ASP.NET"},
	{Header: "X-AspNetMvc-Version", Category: "Framework", Name: "ASP.NET MVC"},
	{Header: "X-Generator", Category: "CMS"},
	{Header: "X-Drupal-Cache", Category: "CMS", Name: "Drupal"},
	{Header: "X-Varnish", Category: "Cache", Name: "Varnish"},
	{Header: "X-Nginx-Cache-Status", Category: "Web Server", Name: "nginx"},
}

// cookieFingerprint maps a cookie name pattern to a technology.
type cookieFingerprint struct {
	Cookie   string
	Name     string
	Category string
}

var cookieFingerprints = []cookieFingerprint{
	{Cookie: "PHPSESSID", Name: "PHP", Category: "Language"},
	{Cookie: "csrftoken", Name: "Django", Category: "Framework"},
	{Cookie: "laravel_session", Name: "Laravel", Category: "Framework"},
	{Cookie: "JSESSIONID", Name: "Java", Category: "Language"},
	{Cookie: "rack.session", Name: "Ruby on Rails", Category: "Framework"},
	{Cookie: "wp-settings", Name: "WordPress", Category: "CMS"},
	{Cookie: "ASP.NET_SessionId", Name: "ASP.NET", Category: "Framework"},
	{Cookie: "_gh_sess", Name: "GitHub", Category: "Platform"},
	{Cookie: "ci_session", Name: "CodeIgniter", Category: "Framework"},
	{Cookie: "SERVERID", Name: "HAProxy", Category: "Load Balancer"},
	{Cookie: "express:sess", Name: "Express.js", Category: "Framework"},
}

// headerValueFingerprint matches specific values within common headers.
type headerValueFingerprint struct {
	Header   string
	Contains string
	Name     string
	Category string
}

var headerValueFingerprints = []headerValueFingerprint{
	// CDN/WAF detection via headers
	{Header: "CF-RAY", Contains: "", Name: "Cloudflare", Category: "CDN"},
	{Header: "X-CDN", Contains: "Incapsula", Name: "Imperva Incapsula", Category: "WAF"},
	{Header: "X-Akamai-Transformed", Contains: "", Name: "Akamai", Category: "CDN"},
	{Header: "X-Cache", Contains: "cloudfront", Name: "AWS CloudFront", Category: "CDN"},
	{Header: "X-Amz-Cf-Id", Contains: "", Name: "AWS CloudFront", Category: "CDN"},
	{Header: "X-Azure-Ref", Contains: "", Name: "Azure CDN", Category: "CDN"},
	{Header: "Via", Contains: "vegur", Name: "Heroku", Category: "PaaS"},
	{Header: "Via", Contains: "varnish", Name: "Varnish", Category: "Cache"},
	{Header: "X-Vercel-Id", Contains: "", Name: "Vercel", Category: "PaaS"},
	{Header: "X-Netlify-Request-ID", Contains: "", Name: "Netlify", Category: "PaaS"},
	{Header: "Fly-Request-Id", Contains: "", Name: "Fly.io", Category: "PaaS"},
	{Header: "X-Shopify-Stage", Contains: "", Name: "Shopify", Category: "E-commerce"},
}

var bodyPatterns = []struct {
	regex    *regexp.Regexp
	name     string
	category string
	version  int // capture group index for version (0 = no version)
}{
	// CMS
	{regex: regexp.MustCompile(`(?i)wp-content/|wp-includes/`), name: "WordPress", category: "CMS"},
	{regex: regexp.MustCompile(`(?i)Drupal\.settings`), name: "Drupal", category: "CMS"},
	{regex: regexp.MustCompile(`(?i)/media/jui/`), name: "Joomla", category: "CMS"},
	// Meta generator
	{regex: regexp.MustCompile(`(?i)<meta[^>]+name=["']generator["'][^>]+content=["']WordPress\s*([\d.]*)`), name: "WordPress", category: "CMS", version: 1},
	{regex: regexp.MustCompile(`(?i)<meta[^>]+name=["']generator["'][^>]+content=["']Drupal\s*([\d.]*)`), name: "Drupal", category: "CMS", version: 1},
	{regex: regexp.MustCompile(`(?i)<meta[^>]+name=["']generator["'][^>]+content=["']Joomla[^"']*([\d.]*)`), name: "Joomla", category: "CMS", version: 1},
	{regex: regexp.MustCompile(`(?i)<meta[^>]+name=["']generator["'][^>]+content=["']Hugo\s*([\d.]*)`), name: "Hugo", category: "Static Site Generator", version: 1},
	{regex: regexp.MustCompile(`(?i)<meta[^>]+name=["']generator["'][^>]+content=["']Jekyll\s*([\d.]*)`), name: "Jekyll", category: "Static Site Generator", version: 1},
	// JS frameworks
	{regex: regexp.MustCompile(`(?i)<div[^>]+id=["']__next["']`), name: "Next.js", category: "JavaScript Framework"},
	{regex: regexp.MustCompile(`(?i)<div[^>]+id=["']__nuxt["']`), name: "Nuxt.js", category: "JavaScript Framework"},
	{regex: regexp.MustCompile(`(?i)<div[^>]+id=["']app["'][^>]*>|<div[^>]+id=["']root["'][^>]*>`), name: "React", category: "JavaScript Framework"},
	{regex: regexp.MustCompile(`(?i)ng-version=["']([\d.]+)["']`), name: "Angular", category: "JavaScript Framework", version: 1},
	{regex: regexp.MustCompile(`(?i)data-v-[a-f0-9]+`), name: "Vue.js", category: "JavaScript Framework"},
	// JS libraries
	{regex: regexp.MustCompile(`(?i)jquery[.-]?([\d.]+)?\.(?:min\.)?js`), name: "jQuery", category: "JavaScript Library", version: 1},
	{regex: regexp.MustCompile(`(?i)bootstrap[.-]?([\d.]+)?\.(?:min\.)?(?:css|js)`), name: "Bootstrap", category: "CSS Framework", version: 1},
	{regex: regexp.MustCompile(`(?i)tailwindcss`), name: "Tailwind CSS", category: "CSS Framework"},
	// Analytics / Tag managers
	{regex: regexp.MustCompile(`(?i)google-analytics\.com/|gtag/js\?id=|googletagmanager\.com`), name: "Google Analytics", category: "Analytics"},
	{regex: regexp.MustCompile(`(?i)connect\.facebook\.net/|fbq\(`), name: "Facebook Pixel", category: "Analytics"},
	// Other
	{regex: regexp.MustCompile(`(?i)cdn\.shopify\.com`), name: "Shopify", category: "E-commerce"},
	{regex: regexp.MustCompile(`(?i)static\.squarespace\.com`), name: "Squarespace", category: "CMS"},
	{regex: regexp.MustCompile(`(?i)cdn\.wix\.com|wixstatic\.com`), name: "Wix", category: "CMS"},
}

// serverVersionPattern extracts version from common server header values.
var serverVersionPattern = regexp.MustCompile(`(?i)^([a-zA-Z][a-zA-Z0-9. _-]*?)(?:[/ ]([\d]+(?:\.[\d]+)*))?$`)

// FingerprintTechnologies detects technologies from HTTP response headers and body.
func FingerprintTechnologies(resp *http.Response, body []byte) []types.Technology {
	seen := make(map[string]bool)
	var techs []types.Technology

	add := func(name, version, category string) {
		key := strings.ToLower(name)
		if seen[key] {
			return
		}
		seen[key] = true
		techs = append(techs, types.Technology{
			Name:     name,
			Version:  version,
			Category: category,
		})
	}

	// 1. Header-based detection
	for _, fp := range headerFingerprints {
		val := resp.Header.Get(fp.Header)
		if val == "" {
			continue
		}
		if fp.Name != "" {
			add(fp.Name, extractVersion(val), fp.Category)
		} else {
			name, version := parseServerHeader(val)
			if name != "" {
				add(name, version, fp.Category)
			}
		}
	}

	// 2. Header value matching (CDN/WAF/PaaS)
	for _, fp := range headerValueFingerprints {
		val := resp.Header.Get(fp.Header)
		if val == "" {
			continue
		}
		if fp.Contains == "" || strings.Contains(strings.ToLower(val), strings.ToLower(fp.Contains)) {
			add(fp.Name, "", fp.Category)
		}
	}

	// 3. Cookie-based detection
	for _, cookie := range resp.Cookies() {
		for _, fp := range cookieFingerprints {
			if strings.EqualFold(cookie.Name, fp.Cookie) {
				add(fp.Name, "", fp.Category)
			}
		}
	}

	// 4. Body-based detection
	if len(body) > 0 {
		bodyStr := string(body)
		for _, bp := range bodyPatterns {
			matches := bp.regex.FindStringSubmatch(bodyStr)
			if matches == nil {
				continue
			}
			version := ""
			if bp.version > 0 && bp.version < len(matches) {
				version = matches[bp.version]
			}
			add(bp.name, version, bp.category)
		}
	}

	return techs
}

// parseServerHeader extracts name and version from a Server header value
// like "nginx/1.24.0" or "Apache/2.4.52 (Ubuntu)".
func parseServerHeader(val string) (string, string) {
	// Handle compound values like "Apache/2.4.52 (Ubuntu)"
	val = strings.TrimSpace(val)
	matches := serverVersionPattern.FindStringSubmatch(val)
	if matches != nil {
		name := strings.TrimSpace(matches[1])
		version := ""
		if len(matches) > 2 {
			version = matches[2]
		}
		return name, version
	}
	// Fallback: just use the raw value as name
	if val != "" {
		return val, ""
	}
	return "", ""
}

// extractVersion tries to pull a version number from a header value.
func extractVersion(val string) string {
	re := regexp.MustCompile(`([\d]+(?:\.[\d]+)+)`)
	if m := re.FindString(val); m != "" {
		return m
	}
	return ""
}
