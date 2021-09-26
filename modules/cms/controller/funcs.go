package controller

import (
	"fmt"
	sysRouter "github.com/iagapie/go-spring/modules/sys/router"
	"github.com/iagapie/go-spring/modules/sys/view"
	"strings"
)

func funcs(ctr *controller) view.FuncMap {
	return view.FuncMap{
		"page":      ctr.renderPage,
		"partial":   ctr.renderPartial,
		"component": ctr.renderComponent,
		"is_page": func(name string) bool {
			return strings.EqualFold(ctr.cur.Page.Name(), name)
		},
		"page_url": func(name string, params ...view.Param) string {
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
	}
}
