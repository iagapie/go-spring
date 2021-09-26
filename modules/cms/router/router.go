package router

import (
	"github.com/iagapie/go-spring/modules/cms/theme"
	"github.com/iagapie/go-spring/modules/sys/router"
)

type Router struct {
	t         theme.Theme
	sysRouter router.Router
	params    router.Params
	url       string
}

func NewRouter(t theme.Theme) *Router {
	return &Router{
		t: t,
	}
}

func (r *Router) Reset() {
	r.t.ResetViews()
	r.sysRouter = nil
}

func (r *Router) URL() string {
	return r.url
}

func (r *Router) Params() router.Params {
	return r.params
}

func (r *Router) SetParams(params router.Params) {
	r.params = params
}

func (r *Router) Param(name string) string {
	if r.params != nil {
		if value, ok := r.params[name]; ok {
			return value
		}
	}
	return ""
}

func (r *Router) FindByURL(url string) theme.View {
	r.url = url
	url = router.NormalizeUrl(url)

	for pass := 1; pass <= 2; pass++ {
		sr := r.getSysRouter()
		if sr.Match(url) {
			r.params = sr.Params()
			name := sr.Matched()

			if page := r.t.Page(name); page != nil && page.Exists() {
				return page
			}

			if pass == 1 {
				r.Reset()
			}
		}
	}
	return nil
}

func (r *Router) FindByPageName(name string, params router.Params) string {
	return r.getSysRouter().URL(name, params)
}

func (r *Router) getSysRouter() router.Router {
	if r.sysRouter != nil {
		return r.sysRouter
	}

	r.sysRouter = router.New()
	for name, page := range r.t.Pages() {
		if pattern := page.Prop("url"); len(pattern) > 0 {
			r.sysRouter.Route(name, pattern)
		}
	}
	r.sysRouter.Sort()

	return r.sysRouter
}
