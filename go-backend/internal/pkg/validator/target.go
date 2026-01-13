package validator

import (
	"net"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Target type constants
const (
	TargetTypeDomain = "domain"
	TargetTypeIP     = "ip"
	TargetTypeCIDR   = "cidr"
)

// IsURLMatchTarget checks if URL hostname matches target
// Returns true if the URL's hostname belongs to the target
//
// Matching rules by target type:
//   - domain: hostname equals target or ends with .target
//   - ip: hostname must exactly equal target
//   - cidr: hostname must be an IP within the CIDR range
func IsURLMatchTarget(urlStr, targetName, targetType string) bool {
	if urlStr == "" || targetName == "" {
		return false
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	hostname := strings.ToLower(parsed.Hostname())
	if hostname == "" {
		return false
	}

	targetName = strings.ToLower(targetName)

	switch targetType {
	case TargetTypeDomain:
		// hostname equals target or ends with .target
		return hostname == targetName || strings.HasSuffix(hostname, "."+targetName)

	case TargetTypeIP:
		// hostname must exactly equal target
		return hostname == targetName

	case TargetTypeCIDR:
		// hostname must be an IP within the CIDR range
		ip := net.ParseIP(hostname)
		if ip == nil {
			return false
		}
		_, network, err := net.ParseCIDR(targetName)
		if err != nil {
			return false
		}
		return network.Contains(ip)

	default:
		return false
	}
}

// IsSubdomainMatchTarget checks if subdomain belongs to target domain
// Returns true if subdomain equals target or ends with .target
func IsSubdomainMatchTarget(subdomain, targetDomain string) bool {
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))
	targetDomain = strings.ToLower(strings.TrimSpace(targetDomain))

	if subdomain == "" || targetDomain == "" {
		return false
	}

	return subdomain == targetDomain || strings.HasSuffix(subdomain, "."+targetDomain)
}

// DetectTargetType auto-detects target type from input string.
// Returns empty string if the input is not a valid target format.
func DetectTargetType(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	// Check CIDR first (must be before IP check since "10.0.0.0" is valid for both)
	if _, _, err := net.ParseCIDR(name); err == nil {
		return TargetTypeCIDR
	}

	// Check IP
	if net.ParseIP(name) != nil {
		return TargetTypeIP
	}

	// Check domain
	if govalidator.IsDNSName(name) {
		return TargetTypeDomain
	}

	return ""
}
