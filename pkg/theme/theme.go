package theme

import (
	"fmt"
	"github.com/iagapie/go-spring/pkg/helper"
	"path/filepath"
	"strings"
)

type (
	Config struct {
		Name         string
		Author       string
		Homepage     string
		Description  string
		PreviewImage string // eg: assets/images/preview.png
	}

	Theme interface {
		ID() string
		BasePath() string
		Dir() string
		SetDir(dir string)
		Cfg() (Config, error)
		Pages() map[string]Page
		Layouts() map[string]Layout
		Partial(name string) Partial
		Reset()
		IsActive() bool
		Datasource() Datasource
		Assets() (uri string, path string)
	}

	theme struct {
		basePath   string
		dir        string
		cfgFile    string
		datasource Datasource
		pageMap    map[string]Page
		layoutMap  map[string]Layout
		partialMap map[string]Partial
	}
)

var (
	activeTheme string
	themesPath  string
)

func SetActiveTheme(theme string) {
	activeTheme = theme
}

func SetThemesPath(path string) {
	themesPath = path
}

func ActiveTheme() Theme {
	if t, ok := Themes()[activeTheme]; ok {
		return t
	}
	return nil
}

func Themes() map[string]Theme {
	all := make(map[string]Theme)
	dirs, err := filepath.Glob(fmt.Sprintf("%s/*", themesPath))
	if err != nil {
		return all
	}
	for _, path := range dirs {
		name := filepath.Base(path)
		all[name] = NewTheme(themesPath, name)
	}
	return all
}

func NewTheme(basePath, dirName string) Theme {
	t := &theme{
		basePath: basePath,
		dir:      dirName,
		cfgFile:  fmt.Sprintf("%s/%s/theme", basePath, dirName),
		partialMap: make(map[string]Partial),
	}
	t.datasource = NewFileDatasource(t)
	return t
}

func (t *theme) ID() string {
	return strings.ToLower(strings.ReplaceAll(t.dir, "/", "-"))
}

func (t *theme) BasePath() string {
	return t.basePath
}

func (t *theme) Dir() string {
	return t.dir
}

func (t *theme) SetDir(dir string) {
	t.dir = dir
}

func (t *theme) Cfg() (Config, error) {
	var cfg Config
	err := helper.ReadConfig(&cfg, t.cfgFile)
	return cfg, err
}

func (t *theme) Pages() map[string]Page {
	if t.pageMap == nil {
		t.pageMap = PagesInTheme(t)
	}
	return t.pageMap
}

func (t *theme) Layouts() map[string]Layout {
	if t.layoutMap == nil {
		t.layoutMap = LayoutsInTheme(t)
	}
	return t.layoutMap
}

func (t *theme) Partial(name string) Partial {
	if p, ok := t.partialMap[name]; ok && p != nil {
		return p
	}
	t.partialMap[name] = PartialInTheme(t, name)
	return t.partialMap[name]
}

func (t *theme) Reset() {
	t.pageMap = nil
	t.layoutMap = nil
	t.partialMap = make(map[string]Partial)
}

func (t *theme) IsActive() bool {
	active := ActiveTheme()
	return active != nil && active.Dir() == t.Dir()
}

func (t *theme) Datasource() Datasource {
	return t.datasource
}

func (t *theme) Assets() (uri string, path string) {
	uri = fmt.Sprintf("/themes/%s/assets", t.dir)
	path = fmt.Sprintf("%s/%s/assets", t.basePath, t.dir)
	return
}
