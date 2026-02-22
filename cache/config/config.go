// Package config provides cache backed by YAML config files only. No database used.
// File naming: dataDir/config/pluginName/PluginName-typeName.yaml
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

const configDirName = "config"

// ConfigCache is a cache backed by a YAML config file. No database.
type ConfigCache struct {
	dataDir    string
	pluginName string
	typeName   string
	path       string
	mu         sync.RWMutex
	cached     any
}

// ConfigFilePath returns the default config file path. Rule: dataDir/config/pluginName/PluginName-typeName.yaml
func ConfigFilePath(dataDir, pluginName, typeName string) string {
	fileName := pluginName + "-" + typeName + ".yaml"
	return filepath.Join(dataDir, configDirName, pluginName, fileName)
}

// ConfigDir returns the config directory for the plugin: dataDir/config/pluginName
func ConfigDir(dataDir, pluginName string) string {
	return filepath.Join(dataDir, configDirName, pluginName)
}

// NewConfigCache creates a config-file cache. typeName is used in the file name: {PluginName}-{typeName}.yaml
func NewConfigCache(dataDir, pluginName, typeName string) *ConfigCache {
	path := ConfigFilePath(dataDir, pluginName, typeName)
	return &ConfigCache{
		dataDir:    dataDir,
		pluginName: pluginName,
		typeName:   typeName,
		path:       path,
	}
}

// Path returns the config file path.
func (c *ConfigCache) Path() string {
	return c.path
}

// Init writes initial config to file when the file does not exist. Call this first with config from initialization (e.g. *types.PluginEngine).
func (c *ConfigCache) Init(initial any) error {
	if c.Exists() {
		return nil
	}
	return c.Save(initial)
}

// Load reads the YAML file into dest and keeps a copy in memory. If the file does not exist, dest is unchanged and no error.
func (c *ConfigCache) Load(dest any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("config cache read: %w", err)
	}
	if len(data) == 0 {
		return nil
	}
	if err := yaml.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("config cache unmarshal: %w", err)
	}
	c.cached = dest
	return nil
}

// Save writes v to the YAML file and updates the in-memory cache.
func (c *ConfigCache) Save(v any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	dir := ConfigDir(c.dataDir, c.pluginName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("config cache mkdir: %w", err)
	}
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("config cache marshal: %w", err)
	}
	if err := os.WriteFile(c.path, data, 0644); err != nil {
		return fmt.Errorf("config cache write: %w", err)
	}
	c.cached = v
	return nil
}

// Get returns the last loaded or saved value from memory. Nil if never Load/Save.
func (c *ConfigCache) Get() any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cached
}

// Exists reports whether the config file exists.
func (c *ConfigCache) Exists() bool {
	_, err := os.Stat(c.path)
	return err == nil
}
