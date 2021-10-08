package config

type Cfg struct {
	App   App   `env-prefix:"APP_" yaml:"app" json:"app"`
	CMS   CMS   `env-prefix:"CMS_" yaml:"cms" json:"cms"`
	CORS  CORS  `env-prefix:"CORS_" yaml:"cors" json:"cors"`
	JWT   JWT   `env-prefix:"JWT_" yaml:"jwt" json:"jwt"`
	DB    DB    `env-prefix:"DB_" yaml:"db" json:"db"`
	Redis Redis `env-prefix:"REDIS_" yaml:"redis" json:"redis"`
}
