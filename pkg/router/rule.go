package router

import (
	"github.com/iagapie/go-spring/pkg/helper"
	"strings"
)

type (
	rule struct {
		name                string
		pattern             string
		staticURL           string
		segments            []string
		staticSegmentCount  int
		dynamicSegmentCount int
		wildSegmentCount    int
	}
)

func newRule(name, pattern string) *rule {
	r := &rule{
		name:     name,
		pattern:  pattern,
		segments: helper.SegmentizeUrl(pattern),
	}
	staticSegments := make([]string, 0)
	for _, segment := range r.segments {
		if strings.HasPrefix(segment, ":") {
			r.dynamicSegmentCount++
			if helper.SegmentIsWildcard(segment) {
				r.wildSegmentCount++
			}
		} else {
			staticSegments = append(staticSegments, segment)
			r.staticSegmentCount++
		}
	}
	r.staticURL = helper.RebuildUrl(staticSegments)
	return r
}

func (r *rule) resolveUrl(url string) (Params, bool) {
	params := make(Params)
	urlSegments := helper.SegmentizeUrl(url)
	var wildSegments []string

	if r.wildSegmentCount == 1 {
		urlSegments, wildSegments = r.captureWildcardSegments(urlSegments)
	}

	if len(urlSegments) > len(r.segments) {
		return params, false
	}

	for index, segment := range r.segments {
		if strings.HasPrefix(segment, ":") {
			paramName := helper.ParameterName(segment)
			params[paramName] = ""

			optional := helper.SegmentIsOptional(segment)

			if optional && index < (len(r.segments)-1) {
				for i := index+1; i < len(r.segments); i++ {
					if !helper.SegmentIsOptional(r.segments[i]) {
						optional = false
						break
					}
				}
			}

			urlSegmentExists := len(urlSegments) > index

			if optional && !urlSegmentExists {
				params[paramName] = helper.SegmentDefaultValue(segment)
				continue
			}

			if !optional && !urlSegmentExists {
				return params, false
			}

			if r := helper.SegmentRegexp(segment); r != nil {
				if !r.MatchString(urlSegments[index]) {
					return params, false
				}
			}

			params[paramName] = urlSegments[index]

			if helper.SegmentIsWildcard(segment) && len(wildSegments) > 0 {
				params[paramName] += helper.RebuildUrl(wildSegments)
			}
		} else if len(urlSegments) <= index || !strings.EqualFold(segment, urlSegments[index]) {
			return params, false
		}
	}

	return params, true
}

func (r *rule) captureWildcardSegments(urlSegments []string) ([]string, []string) {
	newUrlSegments := make([]string, 0)
	wildSegments := make([]string, 0)
	segmentDiff := len(urlSegments) - len(r.segments)
	wildMode := false
	wildCount := 0
	jump := false

	for index, urlSegment := range urlSegments {
		if !jump {
			if wildMode {
				if wildCount < segmentDiff {
					wildSegments = append(wildSegments, urlSegment)
					wildCount++
					continue
				}
				jump = true
			} else if helper.SegmentIsWildcard(r.segments[index]) {
				wildMode = true
			}
		}

		newUrlSegments = append(newUrlSegments, urlSegment)
	}

	return newUrlSegments, wildSegments
}
