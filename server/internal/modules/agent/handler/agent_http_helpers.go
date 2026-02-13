package handler

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	agentdomain "github.com/yyhuni/lunafox/server/internal/modules/agent/domain"
)

func includesHeartbeat(include string) bool {
	if include == "" {
		return false
	}
	for _, part := range strings.Split(include, ",") {
		if strings.EqualFold(strings.TrimSpace(part), "heartbeat") {
			return true
		}
	}
	return false
}

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

func renderInstallScript(tpl *template.Template, data installTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func contextAgent(c *gin.Context) (*agentdomain.Agent, bool) {
	agentVal, ok := c.Get("agent")
	if !ok {
		return nil, false
	}
	agent, ok := agentVal.(*agentdomain.Agent)
	if !ok || agent == nil {
		return nil, false
	}
	return agent, true
}
