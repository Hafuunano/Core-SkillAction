package types

// PluginEngine holds plugin metadata. Convention: each plugin should define a package-level
// variable of this type (e.g. var Plugin or var Meta) with PluginID, PluginName, PluginType,
// and PluginIsDefaultOn set, so plugin code can use it directly (e.g. Plugin.PluginID) without
// reading from a central config.
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

// NewPluginEngine returns a PluginEngine with the given fields. Use in plugin packages for
// one-line definition, e.g. var Plugin = types.NewPluginEngine("id", "name", "skill", true).
func NewPluginEngine(id, name, pluginType string, isDefaultOn bool) PluginEngine {
	return PluginEngine{
		PluginID:          id,
		PluginName:        name,
		PluginType:        pluginType,
		PluginIsDefaultOn: isDefaultOn,
	}
}
