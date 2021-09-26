package plugin

import (
	"fmt"
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/labstack/echo/v4"
	"io/fs"
	"path/filepath"
	"plugin"
	"strings"
)

type Manager struct {
	log     echo.Logger
	infoMap map[string]Info
	codeMap map[string]*info
}

func New(pluginsPath string, log echo.Logger) (*Manager, error) {
	m := &Manager{
		log:     log,
		infoMap: make(map[string]Info),
		codeMap: make(map[string]*info),
	}

	if err := filepath.Walk(pluginsPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".so") {
			if _, err = m.GetByFile(path); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) All() map[string]Info {
	return m.infoMap
}

func (m *Manager) GetByFile(file string) (Info, error) {
	if p, ok := m.infoMap[file]; ok {
		return p, nil
	}

	p, err := plugin.Open(file)
	if err != nil {
		return nil, err
	}

	plug, err := p.Lookup("Plugin")
	if err != nil {
		return nil, err
	}

	i := &info{
		log:  m.log,
		p:    p,
		plug: *plug.(*Plugin),
		file: file,
	}

	details := i.Plugin().Details()
	code := details.Code

	if _, ok := m.codeMap[code]; ok {
		return nil, fmt.Errorf("plugin manager: plugin %s - %s not unique", file, code)
	}

	m.infoMap[file] = i
	m.codeMap[code] = i

	m.log.Debugf("plugin manager: %s[%s] - %s was loaded successfully", details.Name, details.Code, details.Version)

	return m.infoMap[file], nil
}

func (m *Manager) RegisterAll(s *spring.Spring) {
	for _, i := range m.infoMap {
		i.Register(s)
	}
}

func (m *Manager) RoutesAll(f *spring.Frontend, b *spring.Backend) {
	for _, i := range m.infoMap {
		i.Routes(f, b)
	}
}
