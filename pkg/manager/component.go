package manager

import "reflect"

type (
	ComponentDetails struct {
		Code        string `env-required:"required" env:"CODE" yaml:"code" json:"code"`
		Name        string `env-required:"required" env:"NAME" yaml:"name" json:"name"`
		Description string `env:"DESCRIPTION" yaml:"description" json:"description"`
		ViewFile    string `env-default:"default" env:"VIEW_FILE" yaml:"view_file" json:"view_file"`
	}

	ComponentCfgProp struct {
		Title             string       `env-required:"required" env:"TITLE" yaml:"title" json:"title"`
		Description       string       `env:"DESCRIPTION" yaml:"description" json:"description"`
		Default           string       `env:"DEFAULT" yaml:"default" json:"default"`
		Type              reflect.Kind `env-required:"required" env:"TYPE" yaml:"type" json:"type"`
		ValidationPattern string       `env:"VALIDATION_PATTERN" yaml:"validation_pattern" json:"validation_pattern"`
		ValidationMessage string       `env:"VALIDATION_MESSAGE" yaml:"validation_message" json:"validation_message"`
	}

	ComponentCfgProps map[string]ComponentCfgProp
	ComponentProps    map[string]string

	Component interface {
		Details() ComponentDetails
		CfgProps() ComponentCfgProps
		Props() ComponentProps
	}

	ComponentFactory func(props ComponentProps) (Component, error)

	ComponentManager interface {
		Components() map[string]ComponentFactory
		RegisterComponent(fn ComponentFactory, code string, plugin PluginInfo)
		FindPlugin(component Component) PluginInfo
		Has(code string) bool
		Resolve(code string) ComponentFactory
	}

	compManager struct {
		plugManager PluginManager
		components  map[string]ComponentFactory
		plugins     map[string]PluginInfo
	}
)

func NewComponentManager(plugManager PluginManager) ComponentManager {
	return &compManager{
		plugManager: plugManager,
		plugins:     make(map[string]PluginInfo),
	}
}

func (cm *compManager) Components() map[string]ComponentFactory {
	if cm.components == nil {
		cm.load()
	}
	return cm.components
}

func (cm *compManager) RegisterComponent(fn ComponentFactory, code string, plugin PluginInfo) {
	if cm.components == nil {
		cm.components = make(map[string]ComponentFactory)
	}
	cm.components[code] = fn
	if plugin != nil {
		cm.plugins[code] = plugin
	}
}

func (cm *compManager) FindPlugin(component Component) PluginInfo {
	if p, ok := cm.plugins[component.Details().Code]; ok {
		return p
	}
	return nil
}

func (cm *compManager) Has(code string) bool {
	return cm.Resolve(code) != nil
}

func (cm *compManager) Resolve(code string) ComponentFactory {
	if fn, ok := cm.Components()[code]; ok {
		return fn
	}
	return nil
}

func (cm *compManager) load() {
	for _, p := range cm.plugManager.All() {
		for code, fn := range p.RegisterComponents() {
			cm.RegisterComponent(fn, code, p)
		}
	}
}
