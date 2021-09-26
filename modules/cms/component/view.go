package component

import (
	"net/http"
	"reflect"
	"sync"
)

type (
	ViewComponents interface {
		RunComps(r *http.Request) string
		ClearComps()
		AllComps() map[string]Component
		Comp(alias string) Component
		AddComp(alias string, c Component)
		FindByHandler(handler string) Component
	}

	viewComps struct {
		mu    sync.RWMutex
		comps map[string]Component
	}
)

func NewViewComponents() ViewComponents {
	return &viewComps{
		comps: make(map[string]Component),
	}
}

func (v *viewComps) RunComps(r *http.Request) string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	for _, comp := range v.comps {
		if result := comp.OnRun(r); len(result) > 0 {
			return result
		}
	}
	return ""
}

func (v *viewComps) ClearComps() {
	v.comps = make(map[string]Component)
}

func (v *viewComps) AllComps() map[string]Component {
	return v.comps
}

func (v *viewComps) Comp(alias string) Component {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if comp, ok := v.comps[alias]; ok {
		return comp
	}
	return nil
}

func (v *viewComps) AddComp(alias string, c Component) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.comps[alias] = c
}

func (v *viewComps) FindByHandler(handler string) Component {
	v.mu.RLock()
	defer v.mu.RUnlock()
	for _, c := range v.comps {
		if _, ok := reflect.TypeOf(c).MethodByName(handler); ok {
			return c
		}
	}
	return nil
}
