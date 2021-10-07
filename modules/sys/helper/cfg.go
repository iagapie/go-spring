package helper

import (
	"github.com/ilyakaznacheev/cleanenv"
	"strings"
)

var exts = []string{".yml", ".yaml", ".json", ".toml", ".env"}

func ReadConfigWithEnv(cfg interface{}, files ...string) error {
	if err := ReadConfig(cfg, files...); err != nil {
		return err
	}
	return cleanenv.ReadEnv(cfg)
}

func ReadConfig(cfg interface{}, files ...string) error {
	for _, file := range files {
		if hasExt(file) {
			if FileExists(file) {
				if err := cleanenv.ReadConfig(file, cfg); err != nil {
					return err
				}
			}
		} else {
			for _, ext := range exts {
				if FileExists(file + ext) {
					if err := cleanenv.ReadConfig(file+ext, cfg); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func hasExt(file string) bool {
	for _, ext := range exts {
		if strings.HasSuffix(file, ext) {
			return true
		}
	}
	return false
}
