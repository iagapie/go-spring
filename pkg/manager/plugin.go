package manager

import (
	"fmt"
	"github.com/iagapie/go-spring/pkg/spring"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"plugin"
)

type (
	PluginDetails struct {
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
		Details() PluginDetails
	}

	PluginRegister interface {
		Register(*Manager)
	}

	PluginBoot interface {
		Boot(*spring.Spring)
	}

	PluginRoutes interface {
		Routes(*spring.Frontend, *spring.Backend)
	}

	PluginRegisterNavigation interface {
		RegisterNavigation() map[string]MainMenuItem
	}

	PluginRegisterQuickActions interface {
		RegisterQuickActions() map[string]QuickActionItem
	}

	PluginRegisterComponents interface {
		RegisterComponents() map[string]ComponentFactory
	}
)

type (
	PluginManager interface {
		All() map[string]PluginInfo
		GetByFile(file string) (PluginInfo, error)
		Get(code string) PluginInfo
		RegisterAll(*Manager)
		BootAll(*spring.Spring)
		RoutesAll(*spring.Frontend, *spring.Backend)
	}

	PluginInfo interface {
		File() string
		Dir() string
		Plugin() Plugin
		Register(*Manager)
		Boot(*spring.Spring)
		Routes(*spring.Frontend, *spring.Backend)
		RegisterNavigation() map[string]MainMenuItem
		RegisterQuickActions() map[string]QuickActionItem
		RegisterComponents() map[string]ComponentFactory
	}

	plug struct {
		log  *logrus.Entry
		gop  *plugin.Plugin
		p    Plugin
		file string
	}

	plugManager struct {
		log     *logrus.Entry
		plugins map[string]PluginInfo
		codeMap map[string]*plug
	}
)

func NewPluginManager(pluginsPath string, log *logrus.Entry) (PluginManager, error) {
	soFiles, err := filepath.Glob(fmt.Sprintf("%s/*/*.so", pluginsPath))
	if err != nil {
		return nil, err
	}
	pm := &plugManager{
		log:     log,
		plugins: make(map[string]PluginInfo),
		codeMap: make(map[string]*plug),
	}

	for _, soFile := range soFiles {
		if _, err = pm.GetByFile(soFile); err != nil {
			return nil, err
		}
	}

	return pm, nil
}

func (pm *plugManager) All() map[string]PluginInfo {
	return pm.plugins
}

func (pm *plugManager) GetByFile(file string) (PluginInfo, error) {
	if p, ok := pm.plugins[file]; ok {
		return p, nil
	}

	gop, err := plugin.Open(file)
	if err != nil {
		return nil, err
	}

	p, err := gop.Lookup("Plugin")
	if err != nil {
		return nil, err
	}

	info := &plug{
		log:  pm.log,
		gop:  gop,
		p:    *p.(*Plugin),
		file: file,
	}

	details := info.Plugin().Details()
	code := details.Code

	if _, ok := pm.codeMap[code]; ok {
		return nil, fmt.Errorf("plugin manager: plugin %s - %s not unique", file, code)
	}

	pm.plugins[file] = info
	pm.codeMap[code] = info

	pm.log.Debugf("plugin manager: %s[%s] - %s was loaded successfully", details.Name, details.Code, details.Version)

	return pm.plugins[file], nil
}

func (pm *plugManager) Get(code string) PluginInfo {
	if info, ok := pm.codeMap[code]; ok {
		return info
	}
	return nil
}

func (pm *plugManager) RegisterAll(m *Manager) {
	for _, info := range pm.plugins {
		info.Register(m)
	}
}

func (pm *plugManager) BootAll(s *spring.Spring) {
	for _, info := range pm.plugins {
		info.Boot(s)
	}
}

func (pm *plugManager) RoutesAll(f *spring.Frontend, b *spring.Backend) {
	for _, info := range pm.plugins {
		info.Routes(f, b)
	}
}

func (p *plug) File() string {
	return p.file
}

func (p *plug) Dir() string {
	return filepath.Dir(p.file)
}

func (p *plug) Plugin() Plugin {
	return p.p
}

func (p *plug) Register(m *Manager) {
	if r, ok := p.p.(PluginRegister); ok {
		details := p.Plugin().Details()
		p.log.Debugf("plugin info: %s[%s] start register", details.Name, details.Code)
		r.Register(m)
		p.log.Debugf("plugin info: %s[%s] end register", details.Name, details.Code)
	}
}

func (p *plug) Boot(s *spring.Spring) {
	if b, ok := p.p.(PluginBoot); ok {
		details := p.Plugin().Details()
		p.log.Debugf("plugin info: %s[%s] start boot", details.Name, details.Code)
		b.Boot(s)
		p.log.Debugf("plugin info: %s[%s] end boot", details.Name, details.Code)
	}
}

func (p *plug) Routes(f *spring.Frontend, b *spring.Backend) {
	if r, ok := p.p.(PluginRoutes); ok {
		details := p.Plugin().Details()
		p.log.Debugf("plugin info: %s[%s] start routes", details.Name, details.Code)
		r.Routes(f, b)
		p.log.Debugf("plugin info: %s[%s] end routes", details.Name, details.Code)
	}
}

func (p *plug) RegisterNavigation() map[string]MainMenuItem {
	if n, ok := p.p.(PluginRegisterNavigation); ok {
		return n.RegisterNavigation()
	}
	return nil
}

func (p *plug) RegisterQuickActions() map[string]QuickActionItem {
	if n, ok := p.p.(PluginRegisterQuickActions); ok {
		return n.RegisterQuickActions()
	}
	return nil
}

func (p *plug) RegisterComponents() map[string]ComponentFactory {
	if c, ok := p.p.(PluginRegisterComponents); ok {
		return c.RegisterComponents()
	}
	return nil
}
