package view

import "time"

type Param struct {
	Name  string
	Value interface{}
}

var funcs FuncMap

func init() {
	funcs = FuncMap{
		"now": time.Now,
		"now_utc": func() time.Time {
			return time.Now().UTC()
		},
		"format": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
		"param": func(name string, value interface{}) Param {
			return Param{
				Name:  name,
				Value: value,
			}
		},
	}
}

func Add(name string, fn interface{}) {
	funcs[name] = fn
}

func globalFuncs() FuncMap {
	funMap := make(FuncMap)
	for k, v := range funcs {
		funMap[k] = v
	}
	return funMap
}
