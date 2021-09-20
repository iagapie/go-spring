package manager

type Manager struct {
	ConfigManager     ConfigManager
	LogManager        LogManager
	PluginManager     PluginManager
	NavigationManager NavigationManager
	ComponentManager  ComponentManager
}
