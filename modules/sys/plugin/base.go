package plugin

import (
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/labstack/echo/v4"
	"path/filepath"
	"plugin"
)

type (
	Details struct {
		Code        string `env-required:"required" env:"CODE" yaml:"code" json:"code"`
		Name        string `env-required:"required" env:"NAME" yaml:"name" json:"name"`
		Version     string `env-required:"required" env:"VERSION" yaml:"version" json:"version"`
		Description string `env-required:"required" env:"DESCRIPTION" yaml:"description" json:"description"`
		Author      string `env-required:"required" env:"AUTHOR" yaml:"author" json:"author"`
		Icon        string `env:"ICON" yaml:"icon" json:"icon"`
		IconSVG     string `env:"ICON_SVG" yaml:"icon_svg" json:"icon_svg"`
		Homepage    string `env:"HOMEPAGE" yaml:"homepage" json:"homepage"`
	}

	Plugin interface {
		Details() Details
	}

	Register interface {
		Register(s *spring.Spring)
	}

	Routes interface {
		Routes(f *spring.Frontend, b *spring.Backend)
	}

	Info interface {
		File() string
		Dir() string
		Plugin() Plugin
		Register(s *spring.Spring)
		Routes(f *spring.Frontend, b *spring.Backend)
	}

	info struct {
		log  echo.Logger
		p    *plugin.Plugin
		plug Plugin
		file string
	}
)

func (i *info) File() string {
	return i.file
}

func (i *info) Dir() string {
	return filepath.Dir(i.file)
}

func (i *info) Plugin() Plugin {
	return i.plug
}

func (i *info) Register(s *spring.Spring) {
	if r, ok := i.plug.(Register); ok {
		details := i.Plugin().Details()
		i.log.Debugf("plugin info: %s[%s] start register", details.Name, details.Code)
		r.Register(s)
		i.log.Debugf("plugin info: %s[%s] end register", details.Name, details.Code)
	}
}

func (i *info) Routes(f *spring.Frontend, b *spring.Backend) {
	if r, ok := i.plug.(Routes); ok {
		details := i.Plugin().Details()
		i.log.Debugf("plugin info: %s[%s] start routes", details.Name, details.Code)
		r.Routes(f, b)
		i.log.Debugf("plugin info: %s[%s] end routes", details.Name, details.Code)
	}
}
