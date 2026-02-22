// Package cache provides type constants for cache backends: config (no DB) vs database (uses DB).
package cache

const (
	// SourceConfig loads cache from a YAML config file. No database. File naming: {PluginName}-{TypeName}.yaml under data/config/pluginName/.
	SourceConfig = "config"
	// SourceDatabase loads cache from DB and syncs with memory (fastSync). Uses database.
	SourceDatabase = "database"
)
