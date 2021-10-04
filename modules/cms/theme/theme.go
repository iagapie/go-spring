package theme

import (
	"fmt"
	"github.com/iagapie/go-spring/modules/sys/datasource"
	"github.com/iagapie/go-spring/modules/sys/helper"
	"html/template"
	"path/filepath"
)

type (
	Cfg struct {
		Name         string `env-required:"required" env:"THEME_NAME" yaml:"name" json:"name"`
		Author       string `env:"THEME_AUTHOR" yaml:"author" json:"author"`
		Homepage     string `env:"THEME_HOMEPAGE" yaml:"homepage" json:"homepage"`
		Description  string `env:"THEME_DESCRIPTION" yaml:"description" json:"description"`
		PreviewImage string `env-default:"assets/images/preview.png" env:"THEME_PREVIEW_IMAGE" yaml:"preview_image" json:"preview_image"`
	}

	Theme interface {
		IsActive() bool
		Dir() string
		Path() string
		Cfg() (Cfg, error)
		Assets() (uri string, path string)
		Funcs(funcs template.FuncMap)
		ResetViews()
		Pages() ViewMap
		Page(name string) View
		Layout(name string) View
		Partial(name string) View
		ComponentPartial(pluginDir, name string) View
	}

	theme struct {
		basePath   string
		dir        string
		datasource datasource.Datasource
		pages      ViewMap
		layouts    ViewMap
		partials   ViewMap
	}
)

const (
	dPartials   = "partials"
	dLayouts    = "layouts"
	dPages      = "pages"
	dComponents = "components"
	dAssets     = "assets"
)

var (
	_activeTheme string
	_themesPath  string
	_ds          datasource.Datasource
)

func SetThemesPath(themesPath string) {
	_themesPath = themesPath
}

func SetActiveTheme(activeTheme string) {
	_activeTheme = activeTheme
}

func SetDatasource(ds datasource.Datasource) {
	_ds = ds
}

func Themes() map[string]Theme {
	themes := make(map[string]Theme)
	if len(_themesPath) == 0 {
		return themes
	}
	dirs, err := filepath.Glob(fmt.Sprintf("%s/*", _themesPath))
	if err != nil {
		return themes
	}
	for _, p := range dirs {
		name := filepath.Base(p)
		themes[name] = New(_themesPath, name, _ds)
	}
	return themes
}

func ActiveTheme() Theme {
	if len(_activeTheme) == 0 {
		return nil
	}
	return New(_themesPath, _activeTheme, _ds)
}

func New(basePath, dir string, ds datasource.Datasource) Theme {
	t := &theme{
		basePath:   basePath,
		dir:        dir,
		datasource: ds,
	}
	t.ResetViews()
	return t
}

func (t *theme) IsActive() bool {
	return t.Dir() == _activeTheme
}

func (t *theme) Dir() string {
	return t.dir
}

func (t *theme) Path() string {
	return fmt.Sprintf("%s/%s", t.basePath, t.Dir())
}

func (t *theme) Cfg() (Cfg, error) {
	var cfg Cfg
	err := helper.ReadConfig(&cfg, fmt.Sprintf("%s/theme", t.Path()))
	return cfg, err
}

func (t *theme) Assets() (uri string, path string) {
	uri = fmt.Sprintf("/themes/%s/%s", t.dir, dAssets)
	path = fmt.Sprintf("%s/%s", t.Path(), dAssets)
	return
}

func (t *theme) Funcs(funcs template.FuncMap) {
	t.datasource.Funcs(funcs)
}

func (t *theme) ResetViews() {
	t.layouts = make(ViewMap)
	t.pages = make(ViewMap)
	t.partials = make(ViewMap)
}

func (t *theme) Pages() ViewMap {
	if len(t.pages) == 0 {
		for name, v := range t.datasource.Select(fmt.Sprintf("%s/%s", t.Path(), dPages), "html") {
			t.pages[name] = newView(v)
		}
	}
	return t.pages
}

func (t *theme) Page(name string) View {
	if v, ok := t.Pages()[name]; ok {
		return v
	}
	return nil
}

func (t *theme) Layout(name string) View {
	if v, ok := t.layouts[name]; ok {
		return v
	}
	if v := t.datasource.SelectOne(fmt.Sprintf("%s/%s", t.Path(), dLayouts), name, "html"); v != nil {
		t.layouts[name] = newView(v)
		return t.layouts[name]
	}
	return nil
}

func (t *theme) Partial(name string) View {
	if v, ok := t.partials[name]; ok {
		return v
	}
	if v := t.datasource.SelectOne(fmt.Sprintf("%s/%s", t.Path(), dPartials), name, "html"); v != nil {
		t.partials[name] = newView(v)
		return t.partials[name]
	}
	return nil
}

func (t *theme) ComponentPartial(pluginDir, name string) View {
	if p := t.Partial(name); p != nil {
		return p
	}
	if v := t.datasource.SelectOne(fmt.Sprintf("%s/%s", pluginDir, dComponents), name, "html"); v != nil {
		t.partials[name] = newView(v)
		return t.partials[name]
	}
	return nil
}
