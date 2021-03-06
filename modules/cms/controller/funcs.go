package controller

import (
	"fmt"
	sysRouter "github.com/iagapie/go-spring/modules/sys/router"
	"github.com/iagapie/go-spring/modules/sys/view"
	"html/template"
	"strings"
)

func funcs(ctr *controller) template.FuncMap {
	return template.FuncMap{
		"page":      ctr.renderPage,
		"partial":   ctr.renderPartial,
		"component": ctr.renderComponent,
		"isPage": func(name string) bool {
			return strings.EqualFold(ctr.cur.Page.Name(), name)
		},
		"pageURL": func(name string, params ...view.Param) string {
			routerParams := make(sysRouter.Params)
			for _, p := range params {
				routerParams[p.Name] = fmt.Sprintf("%v", p.Value)
			}
			return ctr.router.FindByPageName(name, routerParams)
		},
		"assets": func(name string) string {
			uri, _ := ctr.t.Assets()
			return fmt.Sprintf("%s/%s", uri, name)
		},
		"styles": func() template.HTML {
			assets, _ := ctr.t.Assets()
			styles := ctr.cur.Page.Styles(assets)
			if ctr.cur.Layout != nil {
				styles = fmt.Sprintf("%s%s", ctr.cur.Layout.Styles(assets), styles)
			}
			return template.HTML(styles)
		},
		"scripts": func() template.HTML {
			assets, _ := ctr.t.Assets()
			scripts := ctr.cur.Page.Scripts(assets)
			if ctr.cur.Layout != nil {
				scripts = fmt.Sprintf("%s%s", ctr.cur.Layout.Scripts(assets), scripts)
			}
			return template.HTML(scripts)
		},
	}
}
