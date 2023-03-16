package graph

import xmlparser "webxml/parser"

type FilterDef struct {
	FilterName       string
	FilterClass      string
	Description      string
	InitParamsKeys   []string
	InitParamsValues []string
}

func NewFilterDef(filter xmlparser.Filter) *FilterDef {
	filterDef := &FilterDef{
		FilterName:  filter.FilterName,
		FilterClass: filter.FilterClass,
		Description: filter.Description,
	}

	paramKeys := make([]string, 0)
	paramValues := make([]string, 0)
	for _, param := range filter.InitParams {
		paramKeys = append(paramKeys, param.ParamName)
		paramValues = append(paramValues, param.ParamValue)
	}
	filterDef.InitParamsKeys = paramKeys
	filterDef.InitParamsValues = paramValues

	return filterDef
}
