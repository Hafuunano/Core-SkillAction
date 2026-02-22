// Package config provides YAML config file helpers. No database, no cache struct.
// File naming: dataDir/config/pluginName/PluginName-typeName.yaml
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Hafuunano/Core-SkillAction/types"
	"gopkg.in/yaml.v3"
)

const configDirName = "config"

var (
	mu             sync.RWMutex
	defaultEngine  types.Engine
	defaultPlugin  types.PluginEngine
)

// Set stores engine and plugin for use by NewPathsFromEngine(). Call once at startup.
func Set(engine types.Engine, plugin types.PluginEngine) {
	defaultEngine = engine
	defaultPlugin = plugin
}

func configFilePath(dataDir, pluginName, typeName string) string {
	fileName := pluginName + "-" + typeName + ".yaml"
	return filepath.Join(dataDir, configDirName, pluginName, fileName)
}

func configFilePathFromPlugin(dataDir string, plugin types.PluginEngine, typeName string) string {
	return configFilePath(dataDir, plugin.PluginName, typeName)
}

// Paths holds dataDir and plugin; use NewPathsFromEngine() to get one (after Set).
type Paths struct {
	DataDir string
	Plugin  types.PluginEngine
}

// NewPathsFromEngine returns Paths from the engine and plugin set by Set(). Call Set(engine, plugin) once at startup.
func NewPathsFromEngine() *Paths {
	return &Paths{DataDir: defaultEngine.BotConfigPath, Plugin: defaultPlugin}
}

// Path returns the config file path for typeName (e.g. "engine", "plugin_engine").
func (p *Paths) Path(typeName string) string {
	return configFilePathFromPlugin(p.DataDir, p.Plugin, typeName)
}

// Init writes initial config to file when the file does not exist. Use a default from types (e.g. types.DefaultEngine()).
func Init(path string, initial any) error {
	if Exists(path) {
		return nil
	}
	return Save(path, initial)
}

// Load reads the YAML file into dest. If the file does not exist, dest is unchanged and no error.
func Load(path string, dest any) error {
	mu.RLock()
	defer mu.RUnlock()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("config load: %w", err)
	}
	if len(data) == 0 {
		return nil
	}
	if err := yaml.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("config unmarshal: %w", err)
	}
	return nil
}

// Save writes v to the YAML file. Parent directory is created if needed.
func Save(path string, v any) error {
	mu.Lock()
	defer mu.Unlock()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("config mkdir: %w", err)
	}
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("config marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("config write: %w", err)
	}
	return nil
}

// Exists reports whether the config file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
