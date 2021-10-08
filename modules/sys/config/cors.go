package config

type CORS struct {
	AllowCredentials bool     `env-default:"true" env:"ALLOW_CREDENTIALS" yaml:"allow_credentials" json:"allow_credentials"`
	AllowOrigins     []string `env-default:"*" env:"ALLOW_ORIGINS" yaml:"allow_origins" json:"allow_origins"`
	AllowMethods     []string `env-default:"GET,HEAD,PUT,PATCH,POST,DELETE" env:"ALLOW_METHODS" yaml:"allow_methods" json:"allow_methods"`
	AllowHeaders     []string `env-default:"*" env:"ALLOW_HEADERS" yaml:"allow_headers" json:"allow_headers"`
	ExposeHeaders    []string `env-default:"x-csrf-token,location" env:"EXPOSE_HEADERS" yaml:"expose_headers" json:"expose_headers"`
	MaxAge           int      `env-default:"0" env:"MAX_AGE" yaml:"max_age" json:"max_age"`
}
