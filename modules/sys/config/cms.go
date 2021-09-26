package config

type CMS struct {
	ActiveTheme string `env:"ACTIVE_THEME" env-default:"demo" yaml:"active_theme" json:"active_theme"`
	BackendURI  string `env:"BACKEND_URI" env-default:"/backend" yaml:"backend_uri" json:"backend_uri"`
	PluginsPath string `env:"PLUGINS_PATH" env-default:"./plugins" yaml:"plugins_path" json:"plugins_path"`
	ThemesPath  string `env:"THEMES_PATH" env-default:"./themes" yaml:"themes_path" json:"themes_path"`
}
