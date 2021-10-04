package view

import (
	"github.com/Masterminds/sprig"
	"html/template"
	"time"
)

type Param struct {
	Name  string
	Value interface{}
}

var funcs template.FuncMap

func init() {
	funcs = sprig.FuncMap()
	funcs["nowUtc"] = func() time.Time {
		return time.Now().UTC()
	}
	funcs["param"] = func(name string, value interface{}) Param {
		return Param{
			Name:  name,
			Value: value,
		}
	}
}

func Add(name string, fn interface{}) {
	funcs[name] = fn
}

func globalFuncs() template.FuncMap {
	funMap := make(template.FuncMap)
	for k, v := range funcs {
		funMap[k] = v
	}
	return funMap
}
