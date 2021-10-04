package navigation

import (
	"github.com/iagapie/go-spring/modules/sys/plugin"
	"sync"
)

type Manager struct {
	mui           sync.RWMutex
	muq           sync.RWMutex
	pluginManager *plugin.Manager
	items         map[string]*MainMenuItem
	quickActions  map[string]*QuickActionItem
}

func New(pluginManager *plugin.Manager) *Manager {
	m := &Manager{
		pluginManager: pluginManager,
		items:         make(map[string]*MainMenuItem),
		quickActions:  make(map[string]*QuickActionItem),
	}
	m.LoadMenuItems()
	m.LoadQuickActions()
	return m
}

func (m *Manager) LoadMenuItems() {}

func (m *Manager) LoadQuickActions() {}
