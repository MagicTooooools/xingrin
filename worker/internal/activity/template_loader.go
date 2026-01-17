package activity

import (
	"embed"
	"fmt"
	"sync"
	"text/template"

	"github.com/orbit/worker/internal/pkg"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// TemplateLoader loads and caches command templates from embedded YAML
type TemplateLoader struct {
	fs       embed.FS
	filename string
	once     sync.Once
	cache    map[string]CommandTemplate
	err      error
}

// NewTemplateLoader creates a new template loader
func NewTemplateLoader(fs embed.FS, filename string) *TemplateLoader {
	return &TemplateLoader{
		fs:       fs,
		filename: filename,
	}
}

// Load loads templates (cached with sync.Once)
func (l *TemplateLoader) Load() (map[string]CommandTemplate, error) {
	l.once.Do(func() {
		data, err := l.fs.ReadFile(l.filename)
		if err != nil {
			l.err = fmt.Errorf("failed to read %s: %w", l.filename, err)
			pkg.Logger.Error("Failed to load templates",
				zap.String("file", l.filename),
				zap.Error(l.err))
			return
		}

		l.cache = make(map[string]CommandTemplate)
		if err := yaml.Unmarshal(data, &l.cache); err != nil {
			l.err = fmt.Errorf("failed to parse %s: %w", l.filename, err)
			pkg.Logger.Error("Failed to parse templates",
				zap.String("file", l.filename),
				zap.Error(l.err))
			return
		}

		if err := l.validate(); err != nil {
			l.err = err
			pkg.Logger.Error("Failed to validate templates",
				zap.String("file", l.filename),
				zap.Error(l.err))
			return
		}

		pkg.Logger.Info("Templates loaded",
			zap.String("file", l.filename),
			zap.Int("count", len(l.cache)))
	})

	return l.cache, l.err
}

// Get returns a specific template by name
func (l *TemplateLoader) Get(name string) (CommandTemplate, error) {
	templates, err := l.Load()
	if err != nil {
		return CommandTemplate{}, fmt.Errorf("templates not loaded: %w", err)
	}
	tmpl, ok := templates[name]
	if !ok {
		return CommandTemplate{}, fmt.Errorf("template not found: %s", name)
	}
	return tmpl, nil
}

// validate checks all templates for syntax errors
func (l *TemplateLoader) validate() error {
	for name, tmpl := range l.cache {
		if tmpl.Base == "" {
			return fmt.Errorf("template %s: base command is required", name)
		}
		if _, err := template.New(name).Parse(tmpl.Base); err != nil {
			return fmt.Errorf("template %s: invalid base syntax: %w", name, err)
		}
	}
	return nil
}
