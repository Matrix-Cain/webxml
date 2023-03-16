package graph

import xmlparser "webxml/parser"

type Wrapper struct {
	ServletName      string
	ServletClass     string
	UrlPattern       string
	UrlPatterns      []string
	InitParamsKeys   []string
	InitParamsValues []string
	JSPFile          string
}

func NewWrapper(servlet xmlparser.Servlet, url []string) *Wrapper {
	wrapper := &Wrapper{
		ServletName:  servlet.ServletName,
		ServletClass: servlet.ServletClass,
		UrlPatterns:  url,
		JSPFile:      servlet.JSPFile,
	}

	paramKeys := make([]string, 0)
	paramValues := make([]string, 0)
	for _, param := range servlet.InitParams {
		paramKeys = append(paramKeys, param.ParamName)
		paramValues = append(paramValues, param.ParamValue)
	}
	wrapper.InitParamsKeys = paramKeys
	wrapper.InitParamsValues = paramValues

	return wrapper
}
