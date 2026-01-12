package validator

import (
	"net"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Target type constants
const (
	TargetTypeDomain = "domain"
	TargetTypeIP     = "ip"
	TargetTypeCIDR   = "cidr"
)

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
