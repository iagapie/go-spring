package theme

type (
	Page interface {
		View
		Title() string
		URL() string
		Layout() string
		Description() string
		MetaTitle() string
		MetaDescription() string
		IsHidden() bool
	}

	page struct {
		View
		t Theme
	}
)

const (
	dTitle           = "title"
	dURL             = "url"
	dLayout          = "layout"
	dDescription     = "description"
	dMetaTitle       = "meta_title"
	dMetaDescription = "meta_description"
	dIsHidden        = "is_hidden"
)

func PagesInTheme(t Theme) map[string]Page {
	views := t.Datasource().Select("pages", "html")
	pages := make(map[string]Page)
	for _, v := range views {
		p := &page{
			View: v,
			t:    t,
		}
		pages[p.Name()] = p
	}
	return pages
}

func (p *page) Title() string {
	return p.cfg(dTitle, "")
}

func (p *page) URL() string {
	return p.cfg(dURL, "")
}

func (p *page) Layout() string {
	return p.cfg(dLayout, "default")
}

func (p *page) Description() string {
	return p.cfg(dDescription, "")
}

func (p *page) MetaTitle() string {
	return p.cfg(dMetaTitle, "")
}

func (p *page) MetaDescription() string {
	return p.cfg(dMetaDescription, "")
}

func (p *page) IsHidden() bool {
	return p.Cfg().Section("").Key(dIsHidden).MustInt(0) == 1
}

func (p *page) cfg(key, defaultVal string) string {
	return p.Cfg().Section("").Key(key).MustString(defaultVal)
}
