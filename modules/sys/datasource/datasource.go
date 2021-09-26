package datasource

import (
	"fmt"
	"github.com/iagapie/go-spring/modules/sys/view"
	"github.com/labstack/echo/v4"
	"path/filepath"
	"sync"
)

type (
	Datasource interface {
		Funcs(funcMap view.FuncMap)
		SelectOne(dir, name, ext string) view.View
		Select(dir, ext string) ViewMap
	}

	ViewMap map[string]view.View

	fileDatasource struct {
		mu    sync.RWMutex
		log   echo.Logger
		funcs view.FuncMap
	}
)

func NewFile(log echo.Logger) Datasource {
	return &fileDatasource{
		log:   log,
		funcs: make(view.FuncMap),
	}
}

func (ds *fileDatasource) Funcs(funcMap view.FuncMap) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for name, fn := range funcMap {
		ds.funcs[name] = fn
	}
}

func (ds *fileDatasource) SelectOne(dir, name, ext string) view.View {
	file := fmt.Sprintf("%s/%s.%s", dir, name, ext)
	ds.mu.RLock()
	v := view.New(file, view.WithFuncs(ds.funcs))
	ds.mu.RUnlock()
	if err := v.Load(); err != nil {
		ds.log.Warn(err)
		return nil
	}
	return v
}

func (ds *fileDatasource) Select(dir, ext string) ViewMap {
	views := make(ViewMap)

	files, err := filepath.Glob(fmt.Sprintf("%s/*.%s", dir, ext))
	if err != nil {
		return views
	}

	ds.mu.RLock()
	for _, file := range files {
		v := view.New(file, view.WithFuncs(ds.funcs))
		if err = v.Load(); err != nil {
			ds.log.Warn(err)
			continue
		}
		views[v.Name()] = v
	}
	ds.mu.RUnlock()
	return views
}
