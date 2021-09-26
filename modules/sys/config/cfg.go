package config

type Cfg struct {
	App App `env-prefix:"APP_" yaml:"app" json:"app"`
	CMS CMS `env-prefix:"CMS_" yaml:"cms" json:"cms"`
}
