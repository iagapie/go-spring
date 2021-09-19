package theme

type (
	Layout interface {
		View
	}

	layout struct {
		View
		t Theme
	}
)

func LayoutsInTheme(t Theme) map[string]Layout {
	views := t.Datasource().Select("layouts", "html")
	layouts := make(map[string]Layout)
	for _, v := range views {
		l := &layout{
			View: v,
			t:    t,
		}
		layouts[l.Name()] = l
	}
	return layouts
}
