package helper

import (
	"regexp"
	"strings"
)

func NormalizeUrl(url string) string {
	url = strings.Trim(url, "/")
	return "/" + url
}

func SegmentizeUrl(url string) []string {
	url = NormalizeUrl(url)
	result := make([]string, 0)
	for _, segment := range strings.Split(url, "/") {
		if segment != "" {
			result = append(result, segment)
		}
	}
	return result
}

func RebuildUrl(segments []string) string {
	b := new(strings.Builder)
	for _, segment := range segments {
		if segment != "" {
			b.WriteString("/")
			b.WriteString(segment)
		}
	}
	return NormalizeUrl(b.String())
}

func SegmentDefaultValue(segment string) string {
	optMarkerPos := strings.Index(segment, "?")
	if optMarkerPos == -1 {
		return ""
	}
	value := segment[optMarkerPos+1:]
	if regexMarkerPos := strings.Index(value, "|"); regexMarkerPos != -1 {
		return value[:regexMarkerPos]
	}
	if wildMarkerPos := strings.Index(value, "*"); wildMarkerPos != -1 {
		return value[:wildMarkerPos]
	}
	return value
}

func SegmentIsWildcard(segment string) bool {
	return strings.HasPrefix(segment, ":") && strings.HasSuffix(segment, "*")
}

func ParameterName(segment string) string {
	name := segment[1:]

	optMarkerPos := strings.Index(name, "?")
	wildMarkerPos := strings.Index(name, "*")
	regexMarkerPos := strings.Index(name, "|")

	if wildMarkerPos != -1 {
		if optMarkerPos != -1 {
			return name[:optMarkerPos]
		}
		return name[:wildMarkerPos]
	}

	if optMarkerPos != -1 && regexMarkerPos != -1 {
		if optMarkerPos < regexMarkerPos {
			return name[:optMarkerPos]
		}
		return name[:regexMarkerPos]
	}

	if optMarkerPos != -1 {
		return name[:optMarkerPos]
	}

	if regexMarkerPos != -1 {
		return name[:regexMarkerPos]
	}

	return name
}

func SegmentIsOptional(segment string) bool {
	name := segment[1:]

	optMarkerPos := strings.Index(name, "?")
	if optMarkerPos == -1 {
		return false
	}

	regexMarkerPos := strings.Index(name, "|")
	if regexMarkerPos == -1 {
		return true
	}

	return optMarkerPos < regexMarkerPos
}

func SegmentRegexp(segment string) *regexp.Regexp {
	if pos := strings.Index(segment, "|"); pos != -1 {
		if r, err := regexp.Compile(segment[pos+1:]); err == nil {
			return r
		}
	}
	return nil
}
