package config

type App struct {
	Debug    bool   `env:"DEBUG" env-default:"false" yaml:"debug" json:"debug"`
	Name     string `env:"NAME" env-default:"Spring CMS" yaml:"name" json:"name"`
	Port     int    `env:"PORT" env-default:"80" yaml:"port" json:"port"`
	Timezone string `env:"TIMEZONE" env-default:"UTC" yaml:"timezone" json:"timezone"`
	Locale   string `env:"LOCALE" env-default:"en" yaml:"locale" json:"locale"`
}
