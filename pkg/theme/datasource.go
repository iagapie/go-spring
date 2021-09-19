package theme

import (
	"fmt"
	"github.com/iagapie/go-spring/pkg/helper"
	"path/filepath"
)

type (
	Datasource interface {
		SelectOne(dir, file, ext string) View
		Select(dir, ext string) []View
	}

	fileDatasource struct {
		t        Theme
		basePath string
	}
)

func NewFileDatasource(t Theme) Datasource {
	return &fileDatasource{
		t: t,
		basePath: fmt.Sprintf("%s/%s", t.BasePath(), t.Dir()),
	}
}

func (ds *fileDatasource) SelectOne(dir, file, ext string) View {
	f := fmt.Sprintf("%s/%s/%s.%s", ds.basePath, dir, file, ext)
	if helper.FileExists(f) {
		return NewView(WithViewFile(f))
	}
	return nil
}

func (ds *fileDatasource) Select(dir, ext string) []View {
	views := make([]View, 0)

	files, err := filepath.Glob(fmt.Sprintf("%s/%s/*.%s", ds.basePath, dir, ext))
	if err != nil {
		return views
	}

	for _, file := range files {
		views = append(views, NewView(WithViewFile(file)))
	}
	return views
}
