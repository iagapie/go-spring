package router

import (
	"fmt"
	"github.com/iagapie/go-spring/pkg/helper"
	"sort"
	"strings"
)

type (
	Params map[string]string

	Router interface {
		Reset()
		Route(name, route string)
		Match(url string) bool
		Matched() string
		Parameters() Params
		Sort()
		URI(name string, params ...interface{}) string
		URIFromPattern(pattern string, params ...interface{}) string
		URL(name string, params Params) string
		URLFromPattern(pattern string, params Params) string
	}

	router struct {
		rules        []*rule
		matched      *rule
		params       Params
		defaultValue string
	}
)

func New() Router {
	return &router{
		rules:        make([]*rule, 0),
		params:       make(Params),
		defaultValue: "default",
	}
}

func (r *router) Reset() {
	r.rules = make([]*rule, 0)
}

func (r *router) Route(name, route string) {
	r.rules = append(r.rules, newRule(name, route))
}

func (r *router) Match(url string) bool {
	r.matched = nil
	url = helper.NormalizeUrl(url)

	for _, rule := range r.rules {
		if params, ok := rule.resolveUrl(url); ok {
			r.matched = rule
			r.params = params
			return true
		}
	}

	return false
}

func (r *router) Matched() string {
	if r.matched == nil {
		return ""
	}
	return r.matched.name
}

func (r *router) Parameters() Params {
	return r.params
}

func (r *router) Sort() {
	sort.Slice(r.rules, func(i, j int) bool {
		if r.rules[i].staticSegmentCount > r.rules[j].staticSegmentCount {
			return true
		}

		if r.rules[i].staticSegmentCount == r.rules[j].staticSegmentCount {
			if r.rules[i].dynamicSegmentCount < r.rules[j].dynamicSegmentCount {
				return true
			}
		}

		return false
	})
}

func (r *router) URI(name string, params ...interface{}) string {
	for _, rule := range r.rules {
		if rule.name == name {
			return r.URIFromPattern(rule.pattern, params...)
		}
	}
	return ""
}

func (r *router) URIFromPattern(pattern string, params ...interface{}) string {
	return r.fromPattern(pattern, func(segment string) (string, bool) {
		if len(params) > 0 {
			value := fmt.Sprintf("%v", params[0])
			params = params[1:]
			return value, value != "" && value != helper.SegmentDefaultValue(segment)
		}
		return "", false
	})
}

func (r *router) URL(name string, params Params) string {
	for _, rule := range r.rules {
		if rule.name == name {
			return r.URLFromPattern(rule.pattern, params)
		}
	}
	return ""
}

func (r *router) URLFromPattern(pattern string, params Params) string {
	for param, value := range params {
		if strings.HasPrefix(param, ":") {
			normalizedParam := param[1:]
			params[normalizedParam] = value
			delete(params, param)
		}
	}

	return r.fromPattern(pattern, func(segment string) (string, bool) {
		paramName := helper.ParameterName(segment)
		defaultValue := helper.SegmentDefaultValue(segment)
		if value, ok := params[paramName]; ok && value != "" && value != defaultValue {
			return value, true
		}
		return "", false
	})
}

func (r *router) fromPattern(pattern string, fn func(string) (string, bool)) string {
	url := make([]string, 0)
	lastPopulatedIndex := 0

	for index, segment := range helper.SegmentizeUrl(pattern) {
		if strings.HasPrefix(segment, ":") {
			if value, ok := fn(segment); ok {
				url = append(url, value)
			} else if helper.SegmentIsOptional(segment) {
				if defaultValue := helper.SegmentDefaultValue(segment); defaultValue != "" {
					url = append(url, defaultValue)
				} else {
					url = append(url, r.defaultValue)
				}
				continue
			} else {
				url = append(url, r.defaultValue)
			}
		} else {
			url = append(url, segment)
		}

		lastPopulatedIndex = index
	}

	if len(url) >= lastPopulatedIndex+1 {
		url = url[:lastPopulatedIndex+1]
	}

	return helper.RebuildUrl(url)
}
