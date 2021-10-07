package config

import "time"

type (
	JWTKeys struct {
		Public  string `env-default:"cert/id_rsa.pub" env:"PUBLIC" yaml:"public" json:"public"`
		Private string `env-default:"cert/id_rsa" env:"PRIVATE" yaml:"private" json:"private"`
	}

	JWTDuration struct {
		Access  time.Duration `env-default:"30m" env:"ACCESS" yaml:"access" json:"access"`
		Refresh time.Duration `env-default:"8760h" env:"REFRESH" yaml:"refresh" json:"refresh"` // 365 days
	}

	JWT struct {
		ContextKey  string      `env-default:"user" env:"CONTEXT_KEY" yaml:"context_key" json:"context_key"`
		TokenLookup string      `env-default:"header:Authorization" env:"TOKEN_LOOKUP" yaml:"token_lookup" json:"token_lookup"`
		AuthScheme  string      `env-default:"Bearer" env:"AUTH_SCHEME" yaml:"auth_scheme" json:"auth_scheme"`
		SigningKeys JWTKeys     `env-prefix:"SIGNING_KEYS_" yaml:"signing_keys" json:"signing_keys"`
		TTL         JWTDuration `env-prefix:"TTL_" yaml:"ttl" json:"ttl"`
	}
)
