package manager

import (
	"fmt"
	"github.com/imdario/mergo"
	"sort"
	"strings"
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
		NavItem
		SideMenu map[string]SideMenuItem
	}

	SideMenuItem struct {
		NavItem
		Attributes map[string]interface{}
	}

	QuickActionItem struct {
		NavItem
		Attributes map[string]interface{}
	}
)

type (
	NavigationFunc func(NavigationManager)

	NavigationManager interface {
		RegisterCallback(callback NavigationFunc)
		RegisterMenuItems(owner string, definitions map[string]MainMenuItem)
		AddMainMenuItems(owner string, definitions map[string]MainMenuItem)
		AddMainMenuItem(owner string, code string, definition MainMenuItem)
		GetMainMenuItem(owner string, code string) (MainMenuItem, error)
		RemoveMainMenuItem(owner string, code string)
		AddSideMenuItems(owner string, code string, definitions map[string]SideMenuItem)
		AddSideMenuItem(owner string, code string, sideCode string, definition SideMenuItem) bool
		RemoveSideMenuItems(owner string, code string, sideCodes []string)
		RemoveSideMenuItem(owner string, code string, sideCode string) bool
		RegisterQuickActions(owner string, definitions map[string]QuickActionItem)
		AddQuickActionItems(owner string, definitions map[string]QuickActionItem)
		AddQuickActionItem(owner string, code string, definition QuickActionItem)
		GetQuickActionItem(owner string, code string) (QuickActionItem, error)
		RemoveQuickActionItem(owner string, code string)
		GetMainMenuItems() map[string]MainMenuItem
		GetSideMenuItems(owner string, code string) map[string]SideMenuItem
		GetQuickActionItems() map[string]QuickActionItem
	}

	navManager struct {
		pluginManager PluginManager
		callbacks     []NavigationFunc
		items         map[string]MainMenuItem
		quickActions  map[string]QuickActionItem
		once          sync.Once
	}

	sliceItems []NavItem
)

func NewNavigationManager(pluginManager PluginManager) NavigationManager {
	nm := &navManager{
		pluginManager: pluginManager,
		callbacks:     make([]NavigationFunc, 0),
		items:         make(map[string]MainMenuItem),
		quickActions:  make(map[string]QuickActionItem),
	}

	return nm
}

func (nm *navManager) RegisterCallback(callback NavigationFunc) {
	nm.callbacks = append(nm.callbacks, callback)
}

func (nm *navManager) RegisterMenuItems(owner string, definitions map[string]MainMenuItem) {
	// TODO validation
	nm.AddMainMenuItems(owner, definitions)
}

func (nm *navManager) AddMainMenuItems(owner string, definitions map[string]MainMenuItem) {
	for code, definition := range definitions {
		nm.AddMainMenuItem(owner, code, definition)
	}
}

func (nm *navManager) AddMainMenuItem(owner string, code string, definition MainMenuItem) {
	itemKey := makeItemKey(owner, code)

	if item, ok := nm.items[itemKey]; ok {
		_ = mergo.Merge(&definition, item)
	}

	definition.Code = code
	definition.Owner = owner

	if definition.Order == -1 {
		definition.Order = (len(nm.items) + 1) * 100
	}

	nm.items[itemKey] = definition

	if len(definition.SideMenu) > 0 {
		nm.AddSideMenuItems(owner, code, definition.SideMenu)
	}
}

func (nm *navManager) GetMainMenuItem(owner string, code string) (MainMenuItem, error) {
	itemKey := makeItemKey(owner, code)
	if item, ok := nm.items[itemKey]; ok {
		return item, nil
	}
	return MainMenuItem{}, fmt.Errorf("no main menu item found with key %s", itemKey)
}

func (nm *navManager) RemoveMainMenuItem(owner string, code string) {
	itemKey := makeItemKey(owner, code)
	delete(nm.items, itemKey)
}

func (nm *navManager) AddSideMenuItems(owner string, code string, definitions map[string]SideMenuItem) {
	for sideCode, definition := range definitions {
		nm.AddSideMenuItem(owner, code, sideCode, definition)
	}
}

func (nm *navManager) AddSideMenuItem(owner string, code string, sideCode string, definition SideMenuItem) bool {
	itemKey := makeItemKey(owner, code)
	mainItem, ok := nm.items[itemKey]
	if !ok {
		return false
	}

	definition.Code = sideCode
	definition.Owner = owner

	if item, ok := mainItem.SideMenu[sideCode]; ok {
		_ = mergo.Merge(&definition, item)
	}

	if definition.Order == -1 {
		definition.Order = (len(nm.items[itemKey].SideMenu) + 1) * 100
	}

	nm.items[itemKey].AddSideMenuItem(definition)
	return true
}

func (nm *navManager) RemoveSideMenuItems(owner string, code string, sideCodes []string) {
	for _, sideCode := range sideCodes {
		nm.RemoveSideMenuItem(owner, code, sideCode)
	}
}

func (nm *navManager) RemoveSideMenuItem(owner string, code string, sideCode string) bool {
	itemKey := makeItemKey(owner, code)
	if _, ok := nm.items[itemKey]; ok {
		nm.items[itemKey].RemoveSideMenuItem(sideCode)
		return true
	}
	return false
}

func (nm *navManager) RegisterQuickActions(owner string, definitions map[string]QuickActionItem) {
	// TODO validation
	nm.AddQuickActionItems(owner, definitions)
}

func (nm *navManager) AddQuickActionItems(owner string, definitions map[string]QuickActionItem) {
	for code, definition := range definitions {
		nm.AddQuickActionItem(owner, code, definition)
	}
}

func (nm *navManager) AddQuickActionItem(owner string, code string, definition QuickActionItem) {
	itemKey := makeItemKey(owner, code)

	if item, ok := nm.quickActions[itemKey]; ok {
		_ = mergo.Merge(&definition, item)
	}

	definition.Code = code
	definition.Owner = owner

	if definition.Order == -1 {
		definition.Order = (len(nm.quickActions) + 1) * 100
	}

	nm.quickActions[itemKey] = definition
}

func (nm *navManager) GetQuickActionItem(owner string, code string) (QuickActionItem, error) {
	itemKey := makeItemKey(owner, code)
	if item, ok := nm.quickActions[itemKey]; ok {
		return item, nil
	}
	return QuickActionItem{}, fmt.Errorf("no quick action item found with key %s", itemKey)
}

func (nm *navManager) RemoveQuickActionItem(owner string, code string) {
	itemKey := makeItemKey(owner, code)
	delete(nm.quickActions, itemKey)
}

func (nm *navManager) GetMainMenuItems() map[string]MainMenuItem {
	nm.loadItems()
	return nm.items
}

func (nm *navManager) GetSideMenuItems(owner string, code string) map[string]SideMenuItem {
	if item, ok := nm.GetMainMenuItems()[makeItemKey(owner, code)]; ok {
		return item.SideMenu
	}
	return nil
}

func (nm *navManager) GetQuickActionItems() map[string]QuickActionItem {
	nm.loadItems()
	return nm.quickActions
}

func (nm *navManager) loadItems() {
	nm.once.Do(func() {
		for _, f := range nm.callbacks {
			f(nm)
		}

		for id, plugInfo := range nm.pluginManager.All() {
			if items := plugInfo.RegisterNavigation(); len(items) > 0 {
				nm.RegisterMenuItems(id, items)
			}

			if quickActions := plugInfo.RegisterQuickActions(); len(quickActions) > 0 {
				nm.RegisterQuickActions(id, quickActions)
			}
		}

		nm.sortItems()
		nm.sortQuickActions()

		// TODO filter items by permissions
		// TODO filter quick actions by permissions

		for itemCode := range nm.items {
			if len(nm.items[itemCode].SideMenu) > 0 {
				nm.items[itemCode].sortSide()
				// TODO filter side by permissions
			}
		}
	})
}

func (nm *navManager) sortItems() {
	s := make(sliceItems, 0, len(nm.items))
	for _, item := range nm.items {
		s = append(s, item.NavItem)
	}
	sort.Sort(s)

	items := make(map[string]MainMenuItem)
	for _, item := range s {
		itemKey := makeItemKey(item.Owner, item.Code)
		items[itemKey] = nm.items[itemKey]
	}
	nm.items = items
}

func (nm *navManager) sortQuickActions() {
	s := make(sliceItems, 0, len(nm.quickActions))
	for _, item := range nm.quickActions {
		s = append(s, item.NavItem)
	}
	sort.Sort(s)

	items := make(map[string]QuickActionItem)
	for _, item := range s {
		itemKey := makeItemKey(item.Owner, item.Code)
		items[itemKey] = nm.quickActions[itemKey]
	}
	nm.quickActions = items
}

func (item NavItem) AddPermission(permission string) {
	if item.Permissions == nil {
		item.Permissions = make(map[string]string)
	}
	item.Permissions[permission] = permission
}

func (item NavItem) RemovePermission(permission string) {
	if item.Permissions != nil {
		delete(item.Permissions, permission)
	}
}

func (item NavItem) HasAnyAccess() bool {
	if len(item.Permissions) == 0 {
		return true
	}
	// TODO
	return true
}

func (item MainMenuItem) AddSideMenuItem(sideMenuItem SideMenuItem) {
	if item.SideMenu == nil {
		item.SideMenu = make(map[string]SideMenuItem)
	}
	item.SideMenu[sideMenuItem.Code] = sideMenuItem
}

func (item MainMenuItem) RemoveSideMenuItem(code string) {
	if item.SideMenu != nil {
		delete(item.SideMenu, code)
	}
}

func (item MainMenuItem) GetSideMenuItem(code string) (SideMenuItem, error) {
	if item.SideMenu == nil {
		item.SideMenu = make(map[string]SideMenuItem)
	}

	if side, ok := item.SideMenu[code]; ok {
		return side, nil
	}

	return SideMenuItem{}, fmt.Errorf("no sidenavigation item available with code %s", code)
}

func (item MainMenuItem) sortSide() {
	s := make(sliceItems, 0, len(item.SideMenu))
	for _, i := range item.SideMenu {
		s = append(s, i.NavItem)
	}
	sort.Sort(s)

	items := make(map[string]SideMenuItem)
	for _, i := range s {
		itemKey := makeItemKey(i.Owner, i.Code)
		items[itemKey] = item.SideMenu[itemKey]
	}
	item.SideMenu = items
}

func (item SideMenuItem) AddAttribute(attribute string, value interface{}) {
	if item.Attributes == nil {
		item.Attributes = make(map[string]interface{})
	}
	item.Attributes[attribute] = value
}

func (item SideMenuItem) RemoveAttribute(attribute string) {
	if item.Attributes != nil {
		delete(item.Attributes, attribute)
	}
}

func (item QuickActionItem) AddAttribute(attribute string, value interface{}) {
	if item.Attributes == nil {
		item.Attributes = make(map[string]interface{})
	}
	item.Attributes[attribute] = value
}

func (item QuickActionItem) RemoveAttribute(attribute string) {
	if item.Attributes != nil {
		delete(item.Attributes, attribute)
	}
}

// Len is part of sort.Interface.
func (s sliceItems) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s sliceItems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (s sliceItems) Less(i, j int) bool {
	return s[i].Order < s[j].Order
}

func makeItemKey(owner string, code string) string {
	return strings.ToUpper(fmt.Sprintf("%s.%s", owner, code))
}
