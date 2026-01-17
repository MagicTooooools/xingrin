package activity

// CommandTemplate defines a command template for an activity
type CommandTemplate struct {
	Base     string            `yaml:"base"`     // Base command with required placeholders
	Optional map[string]string `yaml:"optional"` // Optional parameters and their flags
}
