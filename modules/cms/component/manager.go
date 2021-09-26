package component

import (
	"fmt"
	"github.com/iagapie/go-spring/modules/sys/plugin"
	"github.com/iagapie/go-spring/modules/sys/view"
	"sync"
)

type Manager struct {
	mu            sync.RWMutex
	pluginManager *plugin.Manager
	components    map[string]Factory
	infoMap       map[string]plugin.Info
}

func New(pluginManager *plugin.Manager) *Manager {
	m := &Manager{
		pluginManager: pluginManager,
		components:    make(map[string]Factory),
		infoMap:       make(map[string]plugin.Info),
	}
	m.Load()
	return m
}

func (m *Manager) RegisterComponent(fn Factory, code string, info plugin.Info) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.components[code] = fn
	if info != nil {
		m.infoMap[code] = info
	}
}

func (m *Manager) FindPlugin(c Component) plugin.Info {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if info, ok := m.infoMap[c.Details().Code]; ok {
		return info
	}
	return nil
}

func (m *Manager) Has(code string) bool {
	return m.Resolve(code) != nil
}

func (m *Manager) Resolve(code string) Factory {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if fn, ok := m.components[code]; ok {
		return fn
	}
	return nil
}

func (m *Manager) MakeComponent(name string, v view.View, props Props) (Component, error) {
	fn := m.Resolve(name)
	if fn == nil {
		return nil, fmt.Errorf("component factory not found \"%s\", check the component plugin", name)
	}
	return fn(v, props)
}

func (m *Manager) Load() {
	for _, info := range m.pluginManager.All() {
		if regComps, ok := info.Plugin().(PluginRegisterComponents); ok {
			for code, fn := range regComps.RegisterComponents() {
				m.RegisterComponent(fn, code, info)
			}
		}
	}
}
