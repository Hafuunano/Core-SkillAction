package types

type PluginEngine struct {
	PluginID          string `yaml:"plugin_id" json:"plugin_id"`     // Only and only one
	PluginName        string `yaml:"plugin_name" json:"plugin_name"` // what actually called in the code
	PluginType        string `yaml:"plugin_type" json:"plugin_type"`
	PluginIsDefaultOn bool   `yaml:"plugin_is_default_on" json:"plugin_is_default_on"`
}

// DefaultPluginEngine returns a PluginEngine with minimal default values (like Python __init__ defaults).
// Set PluginName (and optionally PluginID, PluginType) when using with config.Init.
func DefaultPluginEngine() PluginEngine {
	return PluginEngine{}
}
