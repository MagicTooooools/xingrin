package subdomain_discovery

import (
	"embed"

	"github.com/orbit/worker/internal/activity"
)

//go:embed templates.yaml
var templatesFS embed.FS

// loader is the template loader for subdomain discovery workflow
var loader = activity.NewTemplateLoader(templatesFS, "templates.yaml")

// getTemplate returns the command template for a given tool
func getTemplate(toolName string) (activity.CommandTemplate, error) {
	return loader.Get(toolName)
}
