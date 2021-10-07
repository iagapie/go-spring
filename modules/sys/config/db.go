package config

type DB struct {
	Host     string `env-default:"localhost" env:"HOST" yaml:"host" json:"host"`
	Port     uint   `env-default:"5432" env:"PORT" yaml:"port" json:"port"`
	Name     string `env-required:"true" env:"NAME" yaml:"name" json:"name"`
	User     string `env-required:"true" env:"USER" yaml:"user" json:"user"`
	Password string `env-required:"true" env:"PASSWORD" yaml:"password" json:"password"`
	SSLMode  string `env-default:"disable" env:"SSL_MODE" yaml:"ssl_mode" json:"ssl_mode"`
	TimeZone string `env-default:"UTC" env:"TIME_ZONE" yaml:"time_zone" json:"time_zone"`
}
