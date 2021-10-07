package config

type Redis struct {
	Addr     string `env-default:"localhost:6379" env:"ADDR" yaml:"addr" json:"addr"`
	Password string `env-default:"" env:"PASSWORD" yaml:"password" json:"password"`
}
