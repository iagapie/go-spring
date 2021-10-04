package navigation

import (
	"sort"
	"sync"
)

type (
	NavItem struct {
		Code        string
		Owner       string
		Label       string
		Icon        string
		IconSVG     string
		URL         string
		Order       int
		Permissions map[string]string
	}

	MainMenuItem struct {
		*NavItem
		sync.RWMutex
		sideMenu map[string]*SideMenuItem
	}

	SideMenuItem struct {
		*NavItem
		Attributes map[string]interface{}
	}

	QuickActionItem struct {
		*NavItem
		Attributes map[string]interface{}
	}

	PluginRegisterNavigation interface {
		RegisterNavigation() map[string]*MainMenuItem
	}

	PluginRegisterQuickActions interface {
		RegisterQuickActions() map[string]*QuickActionItem
	}
)

func NewMainMenuItem(item *NavItem) *MainMenuItem {
	return &MainMenuItem{
		NavItem:  item,
		sideMenu: make(map[string]*SideMenuItem),
	}
}

func (item *MainMenuItem) AddSideMenuItem(sideMenuItem *SideMenuItem) {
	item.Lock()
	defer item.Unlock()
	item.sideMenu[sideMenuItem.Code] = sideMenuItem
}

func (item *MainMenuItem) RemoveSideMenuItem(code string) {
	item.Lock()
	defer item.Unlock()
	delete(item.sideMenu, code)
}

func (item *MainMenuItem) GetSideMenuItem(code string) *SideMenuItem {
	item.RLock()
	defer item.RUnlock()
	if side, ok := item.sideMenu[code]; ok {
		return side
	}
	return nil
}

func (item *MainMenuItem) SideMenu() map[string]*SideMenuItem {
	item.RLock()
	defer item.RUnlock()
	items := make(map[string]*SideMenuItem)
	for k, v := range item.sideMenu {
		items[k] = v
	}
	return items
}

func (item *MainMenuItem) SideMenuCodes() []string {
	item.RLock()
	defer item.RUnlock()
	codes := make([]string, 0, len(item.sideMenu))
	for code, _ := range item.sideMenu {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes
}
