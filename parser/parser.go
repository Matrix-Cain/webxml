package xmlparser

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"os"
)

// WebApp https://docs.oracle.com/cd/E13222_01/wls/docs81/webapp/web_xml.html#1039287
type WebApp struct {
	Filters         []Filter         `xml:"filter,omitempty"`
	FilterMappings  []FilterMapping  `xml:"filter-mapping,omitempty"`
	Listeners       []Listener       `xml:"listener,omitempty"`
	Servlets        []Servlet        `xml:"servlet,omitempty"`
	ServletMappings []ServletMapping `xml:"servlet-mapping,omitempty"`
}

type Filter struct {
	FilterName  string         `xml:"filter-name"`
	FilterClass string         `xml:"filter-class"`
	Description string         `xml:"description,omitempty"`
	InitParams  []ContextParam `xml:"init-param,omitempty"`
}

type FilterMapping struct {
	FilterName      string   `xml:"filter-name"`
	UrlPattern      string   `xml:"url-pattern,omitempty"` // 这里其实有的servlet容器允许多个url-pattern但是目前只采取单映射
	ServletNames    []string `xml:"servlet-name,omitempty"`
	Dispatchers     []string `xml:"dispatcher,omitempty"`
	ServletMappings []string `xml:"servlet-mapping,omitempty"`
}

type Listener struct {
	ListenerClass string `xml:"listener-class"`
}

type Servlet struct {
	ServletName   string         `xml:"servlet-name"`
	ServletClass  string         `xml:"servlet-class"`
	InitParams    []ContextParam `xml:"init-param,omitempty"`
	LoadOnStartup int            `xml:"load-on-startup,omitempty"`
	JSPFile       string         `xml:"jsp-file,omitempty"`
}

type ServletMapping struct {
	ServletName string `xml:"servlet-name"`
	UrlPattern  string `xml:"url-pattern,omitempty"` // 这里其实有的servlet容器允许多个url-pattern但是目前只采取单映射
}

type ContextParam struct {
	ParamName  string `xml:"param-name"`
	ParamValue string `xml:"param-value"`
}

func ParseXML(path string) *WebApp {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var webApp WebApp
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&webApp)
	if err != nil {
		panic(err)
	}

	return &webApp
}
