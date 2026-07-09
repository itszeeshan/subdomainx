package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/types"
)

// WriteTXT writes subdomain results to a text file (one per line)
func WriteTXT(filename string, subdomains []types.SubdomainResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	for _, subdomain := range subdomains {
		if _, err := fmt.Fprintln(file, subdomain.Subdomain); err != nil {
			return err
		}
	}

	return nil
}

// WriteHTTPTXT writes HTTP results to a text file
func WriteHTTPTXT(filename string, httpResults []types.HTTPResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// Write header
	if _, err := fmt.Fprintln(file, "URL\tStatus Code\tTitle\tContent Length\tTechnologies"); err != nil {
		return err
	}

	for _, result := range httpResults {
		technologies := strings.Join(result.Technologies, ",")
		if _, err := fmt.Fprintf(file, "%s\t%d\t%s\t%d\t%s\n",
			result.URL, result.StatusCode, result.Title, result.ContentLength, technologies); err != nil {
			return err
		}
	}

	return nil
}

// WritePortsTXT writes port scan results to a text file
func WritePortsTXT(filename string, portResults []types.PortResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	// Write header
	if _, err := fmt.Fprintln(file, "Host\tIP\tPort\tProtocol\tState\tService\tVersion"); err != nil {
		return err
	}

	for _, result := range portResults {
		for _, port := range result.Ports {
			if _, err := fmt.Fprintf(file, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
				result.Host, result.IP, port.Number, port.Protocol, port.State, port.Service, port.Version); err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteTakeoverTXT writes takeover detection results to a text file.
func WriteTakeoverTXT(filename string, results []types.TakeoverResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if _, err := fmt.Fprintln(file, "Subdomain\tCNAME\tService\tRisk\tEvidence"); err != nil {
		return err
	}

	for _, r := range results {
		if _, err := fmt.Fprintf(file, "%s\t%s\t%s\t%s\t%s\n",
			r.Subdomain, r.CNAME, r.Service, r.Risk, r.Evidence); err != nil {
			return err
		}
	}

	return nil
}

// WriteSubdomainsOnly writes just the subdomain names to a text file
func WriteSubdomainsOnly(filename string, subdomains []types.SubdomainResult) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	for _, subdomain := range subdomains {
		if _, err := fmt.Fprintln(file, subdomain.Subdomain); err != nil {
			return err
		}
	}

	return nil
}
