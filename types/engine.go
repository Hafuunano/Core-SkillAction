package types

// actually they load in memory
type Engine struct {
	BotID             string   `yaml:"bot_id" json:"bot_id"`
	BotName           string   `yaml:"bot_name" json:"bot_name"`
	BotType           string   `yaml:"bot_type" json:"bot_type"`
	BotConfigPath     string   `yaml:"bot_config_path" json:"bot_config_path"`
	BotSuperAdminList []string `yaml:"bot_super_admin_list" json:"bot_super_admin_list"`
}
