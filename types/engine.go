package types

// Engine holds bot config loaded in memory.
type Engine struct {
	BotID             string   `yaml:"bot_id" json:"bot_id"`
	BotName           string   `yaml:"bot_name" json:"bot_name"`
	BotType           string   `yaml:"bot_type" json:"bot_type"`
	BotConfigPath     string   `yaml:"bot_config_path" json:"bot_config_path"`
	BotSuperAdminList []string `yaml:"bot_super_admin_list" json:"bot_super_admin_list"`
}

// DefaultEngine returns an Engine with minimal default values (like Python __init__ defaults).
// Fill in only the fields you need after loading or when passing to config.Init.
func DefaultEngine() Engine {
	return Engine{
		BotSuperAdminList: []string{}, // non-nil so YAML serializes as [] instead of null
	}
}
