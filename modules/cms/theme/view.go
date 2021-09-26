package theme

import (
	"github.com/iagapie/go-spring/modules/cms/component"
	"github.com/iagapie/go-spring/modules/sys/view"
)

type (
	View interface {
		view.View
		component.ViewComponents
	}

	ViewMap map[string]View

	themeView struct {
		view.View
		component.ViewComponents
	}
)

func newView(v view.View) View {
	return &themeView{
		View:           v,
		ViewComponents: component.NewViewComponents(),
	}
}
