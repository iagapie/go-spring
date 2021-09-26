package router

import "strings"

type rule struct {
	name                string
	pattern             string
	staticURL           string
	segments            []string
	staticSegmentCount  int
	dynamicSegmentCount int
	wildSegmentCount    int
}

func newRule(name, pattern string) *rule {
	r := &rule{
		name:     name,
		pattern:  pattern,
		segments: SegmentizeUrl(pattern),
	}
	staticSegments := make([]string, 0)
	for _, segment := range r.segments {
		if strings.HasPrefix(segment, ":") {
			r.dynamicSegmentCount++
			if SegmentIsWildcard(segment) {
				r.wildSegmentCount++
			}
		} else {
			staticSegments = append(staticSegments, segment)
			r.staticSegmentCount++
		}
	}
	r.staticURL = RebuildUrl(staticSegments)
	return r
}

func (r *rule) resolveUrl(url string) (Params, bool) {
	params := make(Params)
	urlSegments := SegmentizeUrl(url)
	var wildSegments []string

	if r.wildSegmentCount == 1 {
		urlSegments, wildSegments = r.captureWildcardSegments(urlSegments)
	}

	if len(urlSegments) > len(r.segments) {
		return params, false
	}

	for index, segment := range r.segments {
		if strings.HasPrefix(segment, ":") {
			paramName := ParameterName(segment)
			params[paramName] = ""

			optional := SegmentIsOptional(segment)

			if optional && index < (len(r.segments)-1) {
				for i := index + 1; i < len(r.segments); i++ {
					if !SegmentIsOptional(r.segments[i]) {
						optional = false
						break
					}
				}
			}

			urlSegmentExists := len(urlSegments) > index

			if optional && !urlSegmentExists {
				params[paramName] = SegmentDefaultValue(segment)
				continue
			}

			if !optional && !urlSegmentExists {
				return params, false
			}

			if re := SegmentRegexp(segment); re != nil {
				if !re.MatchString(urlSegments[index]) {
					return params, false
				}
			}

			params[paramName] = urlSegments[index]

			if SegmentIsWildcard(segment) && len(wildSegments) > 0 {
				params[paramName] += RebuildUrl(wildSegments)
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
			} else if SegmentIsWildcard(r.segments[index]) {
				wildMode = true
			}
		}

		newUrlSegments = append(newUrlSegments, urlSegment)
	}

	return newUrlSegments, wildSegments
}
