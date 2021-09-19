package manager

import (
	"github.com/iagapie/go-spring/pkg/config"
	"github.com/iagapie/go-spring/pkg/helper"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	ConfigManager interface {
		Config() config.Config
		Load() error
	}

	configManger struct {
		files []string
		cfg   *config.Config
	}
)

func NewConfigManager(files ...string) (ConfigManager, error) {
	cm := &configManger{
		files: files,
		cfg:   new(config.Config),
	}

	if err := cm.Load(); err != nil {
		return nil, err
	}

	return cm, nil
}

func (cm *configManger) Config() config.Config {
	return *cm.cfg
}

func (cm *configManger) Load() error {
	if err := helper.ReadConfig(cm.cfg, cm.files...); err != nil {
		return err
	}
	return cleanenv.ReadEnv(cm.cfg)
}
