package helper

import (
	"github.com/ilyakaznacheev/cleanenv"
	"strings"
)

var exts = []string{".yml", ".yaml", ".json", ".toml", ".env"}

func ReadConfig(cfg interface{}, files ...string) (err error) {
	for _, file := range files {
		if hasExt(file) {
			err = readConfig(cfg, file)
		} else {
			for _, ext := range exts {
				err = readConfig(cfg, file+ext)
			}
		}
	}
	return
}

func readConfig(cfg interface{}, file string) error {
	if FileExists(file) {
		if err := cleanenv.ReadConfig(file, cfg); err != nil {
			return err
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
