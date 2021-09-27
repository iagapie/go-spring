package components

import (
	"github.com/iagapie/go-spring/modules/cms/component"
	"github.com/iagapie/go-spring/modules/cms/theme"
	"github.com/iagapie/go-spring/modules/sys/spring"
	"github.com/iagapie/go-spring/modules/sys/view"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"reflect"
)

type todo struct {
	*component.CompBase
	v theme.View
}

func NewTodo(v view.View, props component.Props) component.Component {
	return &todo{
		CompBase: component.NewCompBase(props),
		v:        v.(theme.View),
	}
}

func (*todo) Details() component.Details {
	return component.Details{
		Code:     "todo",
		Name:     "Todo Component",
		ViewFile: "default",
	}
}

func (*todo) CfgProps() component.CfgProps {
	return component.CfgProps{
		"max": {
			Title:   "Max",
			Default: "10",
			Type:    reflect.Int,
		},
		"min": {
			Title:   "Max",
			Default: "0",
			Type:    reflect.Int,
		},
	}
}

func (t *todo) Init(s *spring.Spring) {
	log.Info("component todo Init()")
	log.Info(s.Cfg.App.Name)
	t.v.AddJS("js/test.js", view.Param{
		Name: "defer",
		Value: true,
	})
}

func (*todo) OnRun(r *http.Request) string {
	log.Info("component todo OnRun()")
	log.Info(r.RequestURI)
	return ""
}

type Data struct {
	Foo string `json:"foo"`
}

func (*todo) OnFetchData(c echo.Context) map[interface{}]interface{} {
	var v Data
	c.Bind(&v)
	log.Info(v.Foo)
	return map[interface{}]interface{}{
		"foo": v.Foo,
	}
}

func (*todo) OnFetchForm(c echo.Context) interface{} {
	return map[string]interface{}{"title": c.Request().PostFormValue("title")}
}
