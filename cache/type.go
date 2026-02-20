// Package cache provides type constants and re-exports for cache backends: config (no DB) vs database (uses DB).
package cache

import "github.com/Hafuunano/Core-SkillAction/cache/config"

const (
	// SourceConfig loads cache from a YAML config file. No database. File naming: {PluginName}-{TypeName}.yaml under data/config/pluginName/.
	SourceConfig = "config"
	// SourceDatabase loads cache from DB and syncs with memory (fastSync). Uses database.
	SourceDatabase = "database"
)

// ConfigFilePath returns the default config file path. Rule: dataDir/config/pluginName/PluginName-typeName.yaml
func ConfigFilePath(dataDir, pluginName, typeName string) string {
	return config.ConfigFilePath(dataDir, pluginName, typeName)
}

// ConfigDir returns the config directory for the plugin: dataDir/config/pluginName
func ConfigDir(dataDir, pluginName string) string {
	return config.ConfigDir(dataDir, pluginName)
}
