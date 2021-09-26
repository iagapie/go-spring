package component

import (
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/iagapie/go-spring/modules/sys/view"
	"net/http"
	"reflect"
	"strings"
)

type (
	Details struct {
		Code        string `env-required:"required" env:"CODE" yaml:"code" json:"code"`
		Name        string `env-required:"required" env:"NAME" yaml:"name" json:"name"`
		Description string `env:"DESCRIPTION" yaml:"description" json:"description"`
		ViewFile    string `env-default:"default" env:"VIEW_FILE" yaml:"view_file" json:"view_file"`
	}

	CfgProp struct {
		Title             string       `env-required:"required" env:"TITLE" yaml:"title" json:"title"`
		Description       string       `env:"DESCRIPTION" yaml:"description" json:"description"`
		Default           string       `env:"DEFAULT" yaml:"default" json:"default"`
		Type              reflect.Kind `env-required:"required" env:"TYPE" yaml:"type" json:"type"`
		ValidationPattern string       `env:"VALIDATION_PATTERN" yaml:"validation_pattern" json:"validation_pattern"`
		ValidationMessage string       `env:"VALIDATION_MESSAGE" yaml:"validation_message" json:"validation_message"`
	}

	CfgProps map[string]CfgProp
	Props    map[string]string

	Component interface {
		Details() Details
		CfgProps() CfgProps
		Props() Props
		SetProps(props Props)
		Prop(name string) string
		SetProp(name, value string)
		SetExternalPropName(name, extName string)
		ExternalPropName(name string) string
		ParamName(name string) string
		Alias() string
		SetAlias(alias string)
		Init(s *spring.Spring)
		OnRun(r *http.Request) string
		OnRender() string
	}

	Factory    func(v view.View, props Props) (Component, error)
	FactoryMap map[string]Factory

	PluginRegisterComponents interface {
		RegisterComponents() FactoryMap
	}

	CompBase struct {
		props                 Props
		externalPropertyNames map[string]string
		alias                 string
	}
)

func NewCompBase(props Props) *CompBase {
	return &CompBase{
		props:                 props,
		externalPropertyNames: make(map[string]string),
	}
}

func (comp *CompBase) Props() Props {
	return comp.props
}

func (comp *CompBase) SetProps(props Props) {
	comp.props = props
}

func (comp *CompBase) Prop(name string) string {
	if value, ok := comp.props[name]; ok {
		return value
	}
	return ""
}

func (comp *CompBase) SetProp(name, value string) {
	comp.props[name] = value
}

func (comp *CompBase) SetExternalPropName(name, extName string) {
	comp.externalPropertyNames[name] = extName
}

func (comp *CompBase) ExternalPropName(name string) string {
	if extName, ok := comp.externalPropertyNames[name]; ok {
		return extName
	}
	return ""
}

func (comp *CompBase) ParamName(name string) string {
	if extName := comp.ExternalPropName(name); strings.HasPrefix(extName, ":") {
		return extName[1:]
	}
	return ""
}

func (comp *CompBase) Alias() string {
	return comp.alias
}

func (comp *CompBase) SetAlias(alias string) {
	comp.alias = alias
}

func (comp *CompBase) Init(s *spring.Spring) {
}

func (comp *CompBase) OnRun(r *http.Request) string {
	return ""
}

func (comp *CompBase) OnRender() string {
	return ""
}
