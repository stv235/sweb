package pages

import (
	"net/url"
	"strings"
)

type RouteHandler struct {
	Route        Route
	patternParts []string
}

func (handler RouteHandler) Match(path string) (bool, url.Values) {
	path = strings.Trim(path, "/")

	pathParts := make([]string, 0)

	if path != "" {
		pathParts = strings.Split(path, "/")
	}

	if len(pathParts) != len(handler.patternParts) {
		return false, nil
	}

	parameters := url.Values{}

	for i, patternPart := range handler.patternParts {
		if patternPart[0] == '{' {
			parameters.Set(strings.Trim(patternPart, "{}"), pathParts[i])
		} else if patternPart != pathParts[i] {
			return false, nil
		}
	}

	return true, parameters
}