package graph

import (
	"net/url"
	xmlparser "webxml/parser"
)

type FilterMap struct {
	FilterName      string
	UrlPattern      string
	ServletNames    []string
	Dispatchers     []string
	ServletMappings []string

	matchAllUrlPatterns  bool
	matchAllServletNames bool
}

func NewFilterMap(filterMapping xmlparser.FilterMapping) *FilterMap {
	filterMap := &FilterMap{
		FilterName:      filterMapping.FilterName,
		UrlPattern:      filterMapping.UrlPattern,
		Dispatchers:     filterMapping.Dispatchers,
		ServletMappings: filterMapping.ServletMappings,
	}

	for _, servletName := range filterMapping.ServletNames {
		filterMap.AddServletName(servletName)
	}

	return filterMap
}

func (f *FilterMap) AddURLPattern(urlPattern string) {
	decoded, _ := url.QueryUnescape(urlPattern)
	f.AddURLPatternDecoded(decoded)
}

func (f *FilterMap) AddURLPatternDecoded(urlPattern string) {
	if "*" == urlPattern {
		f.matchAllUrlPatterns = true
	} else {
		decoded, _ := url.QueryUnescape(urlPattern) // 这里双重解码行为和Tomcat8.5.87保持一致
		f.UrlPattern = decoded
	}
}

func (f *FilterMap) AddServletName(servletName string) {
	if "*" == servletName {
		f.matchAllServletNames = true
	} else {
		f.ServletNames = append(f.ServletNames, servletName)
	}
}

func (f *FilterMap) GetMatchAllUrlPatterns() bool {
	return f.matchAllUrlPatterns
}
