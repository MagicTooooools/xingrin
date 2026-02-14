package handler

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func inferServerURL(c *gin.Context) string {
	proto := c.GetHeader("X-Forwarded-Proto")
	if proto == "" {
		if c.Request.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}

	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}
	if host == "" {
		return ""
	}

	port := c.GetHeader("X-Forwarded-Port")
	if port != "" && !isStandardPort(proto, port) {
		host = host + ":" + port
	}
	return fmt.Sprintf("%s://%s", proto, host)
}

func isStandardPort(proto, port string) bool {
	return (proto == "http" && port == "80") || (proto == "https" && port == "443")
}

func getForwardedIP(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if forwarded := strings.TrimSpace(c.GetHeader("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}
	if realIP := strings.TrimSpace(c.GetHeader("X-Real-IP")); realIP != "" {
		return realIP
	}
	return c.ClientIP()
}
