package graph

import (
	"strings"
)

func MatchFiltersURL(filterMap *FilterMap, requestPath string) bool {
	if filterMap.GetMatchAllUrlPatterns() {
		return true
	}

	if requestPath == "" {
		return false
	}

	if MatchFiltersURLWithPath(filterMap.UrlPattern, requestPath) {
		return true
	}

	return false
}

// MatchFiltersURLWithPath mock logic in Tomcat ApplicationFilterFactory logic
func MatchFiltersURLWithPath(testPath string, requestPath string) bool {
	if testPath == "" {
		return false
	}

	// Case 1 - Exact Match
	if testPath == requestPath {
		return true
	}

	// Case 2 - Path Match ("/.../*")
	if testPath == "/*" {
		return true
	}
	if strings.HasSuffix(testPath, "/*") {
		if len(requestPath) >= len(testPath)-2 && testPath[0:len(testPath)-2] == requestPath[0:len(testPath)-2] {
			if testPath[:len(testPath)-3] == requestPath {
				return true
			}
			if '/' == requestPath[len(testPath)-2] {
				return true
			}
		}
		return false
	}

	// Case 3 - Extension Match
	if strings.HasPrefix(testPath, "*.") {
		slash := strings.LastIndex(requestPath, "/")
		period := strings.LastIndex(requestPath, ".")
		if ((slash >= 0) && (period > slash) && (period != len(requestPath)-1)) &&
			(len(requestPath)-period == len(testPath)-1) {
			return period+len(testPath)-1 <= len(requestPath) && testPath[2:] == requestPath[period+1:len(testPath)+period-1]
		}
	}

	// Case 4 - "Default" Match
	return false // NOTE - Not relevant for selecting filters
}
