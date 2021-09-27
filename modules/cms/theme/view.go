package theme

import (
	"fmt"
	"github.com/iagapie/go-spring/modules/cms/component"
	"github.com/iagapie/go-spring/modules/sys/view"
	"html"
	"strings"
	"sync"
)

type (
	View interface {
		view.View
		component.ViewComponents
		AddCSS(name string, attrs ...view.Param)
		AddJS(name string, attrs ...view.Param)
		Styles(assets string) string
		Scripts(assets string) string
	}

	ViewMap map[string]View

	themeView struct {
		view.View
		component.ViewComponents
		mus     sync.RWMutex
		muj     sync.RWMutex
		styles  map[string][]view.Param
		scripts map[string][]view.Param
	}
)

func newView(v view.View) View {
	return &themeView{
		View:           v,
		ViewComponents: component.NewViewComponents(),
		styles:         make(map[string][]view.Param),
		scripts:        make(map[string][]view.Param),
	}
}

func (v *themeView) AddCSS(name string, attrs ...view.Param) {
	v.mus.Lock()
	defer v.mus.Unlock()
	v.styles[name] = attrs
}

func (v *themeView) AddJS(name string, attrs ...view.Param) {
	v.muj.Lock()
	defer v.muj.Unlock()
	v.scripts[name] = attrs
}

func (v *themeView) Styles(assets string) string {
	v.mus.RLock()
	defer v.mus.RUnlock()
	b := new(strings.Builder)
	for name, attrs := range v.styles {
		if !strings.HasPrefix(name, "http") {
			name = fmt.Sprintf("%s/%s", assets, name)
		}
		b.WriteString(fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\"%s>\n", name, toAttrs(attrs)))
	}
	return b.String()
}

func (v *themeView) Scripts(assets string) string {
	v.muj.RLock()
	defer v.muj.RUnlock()
	b := new(strings.Builder)
	for name, attrs := range v.scripts {
		if !strings.HasPrefix(name, "http") {
			name = fmt.Sprintf("%s/%s", assets, name)
		}
		b.WriteString(fmt.Sprintf("<script src=\"%s\"%s></script>\n", name, toAttrs(attrs)))
	}
	return b.String()
}

func toAttrs(data []view.Param) string {
	b := new(strings.Builder)
	for _, attr := range data {
		if v, ok := attr.Value.(bool); ok && v {
			b.WriteString(" " + attr.Name)
			continue
		}
		b.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.Name, html.EscapeString(fmt.Sprintf("%v", attr.Value))))
	}
	return b.String()
}
