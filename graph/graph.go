package graph

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
	"log"
	xmlparser "webxml/parser"
)

var filterMaps = make([]*FilterMap, 0)
var filterDefs = make([]*FilterDef, 0)
var wrappers = make([]*Wrapper, 0)
var filterConfigs = make(map[string]*FilterConfig)

func BuildGraph(neo4jAddr string, neo4jUser string, neo4jPass string, xmlpath string, projectName string, pathVerbose bool) {
	ctx := context.Background()

	driver, err := neo4j.NewDriverWithContext(neo4jAddr, neo4j.BasicAuth(neo4jUser, neo4jPass, ""))
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close(ctx)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(ctx)

	webApp := xmlparser.ParseXML(xmlpath)

	buildMainNodeForApp(projectName, session)

	for _, listener := range webApp.Listeners {
		buildStandAloneNode(listener, session)
	}

	for _, filterDef := range webApp.Filters {
		newDef := NewFilterDef(filterDef)
		filterDefs = append(filterDefs, newDef)
		filterConfigs[filterDef.FilterName] = &FilterConfig{newDef}
	}

	for _, filterMap := range webApp.FilterMappings {
		filterMaps = append(filterMaps, NewFilterMap(filterMap))
	}

	var servletRoute = make(map[string][]string)
	for _, servletMapping := range webApp.ServletMappings {
		servletRoute[servletMapping.ServletName] = append(servletRoute[servletMapping.ServletName], servletMapping.UrlPattern)
	}

	for _, servlet := range webApp.Servlets {
		wrappers = append(wrappers, NewWrapper(servlet, servletRoute[servlet.ServletName]))
	}

	for _, wrapper := range wrappers {
		for _, pattern := range wrapper.UrlPatterns {
			wrapper.UrlPattern = pattern
			buildServletInvocationChain(wrapper, session, pathVerbose)
		}
	}

	fmt.Println("Done")
}

func buildMainNodeForApp(projectName string, session neo4j.SessionWithContext) {
	ctx := context.Background()
	if projectName == "" {
		projectName = "WebApp"
	}

	// set main node representing the web app
	_, err := session.Run(ctx, "MERGE (:App {name: $name})", map[string]interface{}{
		"name": projectName,
	})

	if err != nil {
		panic(err)
	}
}

func buildStandAloneNode(listener *xmlparser.Listener, session neo4j.SessionWithContext) {
	ctx := context.Background()
	_, err := session.Run(ctx, "MERGE (:Listener {class: $class})", map[string]interface{}{
		"class": listener.ListenerClass,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Listener Class: %s\n", listener.ListenerClass)
}

func buildServletInvocationChain(wrapper *Wrapper, session neo4j.SessionWithContext, pathVerbose bool) {
	ctx := context.Background()
	filterChains := make([]*FilterConfig, 0)
	for _, filterMap := range filterMaps {
		if !MatchFiltersURL(filterMap, wrapper.UrlPattern) {
			continue
		}
		filterChains = append(filterChains, filterConfigs[filterMap.FilterName])
	}
	if len(filterChains) > 0 {
		var err error
		for index, filter := range filterChains {
			if index == 0 {
				if pathVerbose {
					_, err = session.Run(ctx, "MATCH (a:App)"+
						"MERGE (b:Filter {name: $name, class: $class, initParamsKeys: $initParamsKeys, initParamsValues: $initParamsValues})"+
						fmt.Sprintf("MERGE (a)-[:`%v` {url: $url}]->(b)", wrapper.ServletName), map[string]interface{}{
						"name":             filter.FilterName,
						"class":            filter.FilterClass,
						"initParamsKeys":   filter.InitParamsKeys,
						"initParamsValues": filter.InitParamsValues,
						"url":              wrapper.UrlPattern,
					})
				} else {
					_, err = session.Run(ctx, "MATCH (a:App)"+
						"MERGE (b:Filter {name: $name, class: $class, initParamsKeys: $initParamsKeys, initParamsValues: $initParamsValues})"+
						"MERGE (a)-[:route]->(b)", map[string]interface{}{
						"name":             filter.FilterName,
						"class":            filter.FilterClass,
						"initParamsKeys":   filter.InitParamsKeys,
						"initParamsValues": filter.InitParamsValues,
						"url":              wrapper.UrlPattern,
					})
				}
				goto ERROR_CHECK
			}
			if pathVerbose {
				_, err = session.Run(ctx, "MATCH (a:Filter {name: $name1, class: $class1, initParamsKeys: $initParamsKeys1, initParamsValues: $initParamsValues1})"+
					"MERGE (b:Filter {name: $name, class: $class, initParamsKeys: $initParamsKeys, initParamsValues: $initParamsValues})"+
					fmt.Sprintf("MERGE (a)-[:`%v` {url: $url}]->(b)", wrapper.ServletName), map[string]interface{}{
					"name":              filter.FilterName,
					"class":             filter.FilterClass,
					"initParamsKeys":    filter.InitParamsKeys,
					"initParamsValues":  filter.InitParamsValues,
					"name1":             filterChains[index-1].FilterName,
					"class1":            filterChains[index-1].FilterClass,
					"initParamsKeys1":   filterChains[index-1].InitParamsKeys,
					"initParamsValues1": filterChains[index-1].InitParamsValues,
					"url":               wrapper.UrlPattern,
				})
			} else {
				_, err = session.Run(ctx, "MATCH (a:Filter {name: $name1, class: $class1, initParamsKeys: $initParamsKeys1, initParamsValues: $initParamsValues1})"+
					"MERGE (b:Filter {name: $name, class: $class, initParamsKeys: $initParamsKeys, initParamsValues: $initParamsValues})"+
					"MERGE (a)-[:route]->(b)", map[string]interface{}{
					"name":              filter.FilterName,
					"class":             filter.FilterClass,
					"initParamsKeys":    filter.InitParamsKeys,
					"initParamsValues":  filter.InitParamsValues,
					"name1":             filterChains[index-1].FilterName,
					"class1":            filterChains[index-1].FilterClass,
					"initParamsKeys1":   filterChains[index-1].InitParamsKeys,
					"initParamsValues1": filterChains[index-1].InitParamsValues,
					"url":               wrapper.UrlPattern,
				})
			}
			goto ERROR_CHECK

		ERROR_CHECK:
			if err != nil {
				panic(err)
			}
			fmt.Printf("Filter  Name: %s, Class: %s\n", filter.FilterName, filter.FilterClass)
		}
		if pathVerbose {
			_, err = session.Run(ctx, "MATCH (a:Filter {name: $name, class: $class, initParamsKeys: $initParamsKeys, initParamsValues: $initParamsValues})"+
				"MERGE (b:Servlet {name: $name1, class: $class1, initParamsKeys: $initParamsKeys1, initParamsValues: $initParamsValues1, jspFile: $jspFile1})"+
				fmt.Sprintf("MERGE (a)-[:`%v` {url: $url}]->(b)", wrapper.ServletName), map[string]interface{}{
				"name":              filterChains[len(filterChains)-1].FilterName,
				"class":             filterChains[len(filterChains)-1].FilterClass,
				"initParamsKeys":    filterChains[len(filterChains)-1].InitParamsKeys,
				"initParamsValues":  filterChains[len(filterChains)-1].InitParamsValues,
				"name1":             wrapper.ServletName,
				"class1":            wrapper.ServletClass,
				"initParamsKeys1":   wrapper.InitParamsKeys,
				"initParamsValues1": wrapper.InitParamsValues,
				"jspFile1":          wrapper.JSPFile,
				"url":               wrapper.UrlPattern,
			})
		} else {
			_, err = session.Run(ctx, "MATCH (a:Filter {name: $name, class: $class, initParamsKeys: $initParamsKeys, initParamsValues: $initParamsValues})"+
				"MERGE (b:Servlet {name: $name1, class: $class1, initParamsKeys: $initParamsKeys1, initParamsValues: $initParamsValues1, jspFile: $jspFile1})"+
				"MERGE (a)-[:route]->(b)", map[string]interface{}{
				"name":              filterChains[len(filterChains)-1].FilterName,
				"class":             filterChains[len(filterChains)-1].FilterClass,
				"initParamsKeys":    filterChains[len(filterChains)-1].InitParamsKeys,
				"initParamsValues":  filterChains[len(filterChains)-1].InitParamsValues,
				"name1":             wrapper.ServletName,
				"class1":            wrapper.ServletClass,
				"initParamsKeys1":   wrapper.InitParamsKeys,
				"initParamsValues1": wrapper.InitParamsValues,
				"jspFile1":          wrapper.JSPFile,
				"url":               wrapper.UrlPattern,
			})

		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Servlet  Name: %s, Class: %s\n", wrapper.ServletName, wrapper.ServletClass)
	} else { // 不存在filter
		var err error
		if pathVerbose {
			_, err = session.Run(ctx, "MATCH (a:App)"+
				"MERGE (b:Servlet {name: $name, class: $class, initParamsKeys: $initParamsKeys, initParamsValues: $initParamsValues, jspFile: $jspFile})"+
				fmt.Sprintf("MERGE (a)-[:`%v` {url: $url}]->(b)", wrapper.ServletName), map[string]interface{}{
				"name":             wrapper.ServletName,
				"class":            wrapper.ServletClass,
				"initParamsKeys":   wrapper.InitParamsKeys,
				"initParamsValues": wrapper.InitParamsValues,
				"jspFile":          wrapper.JSPFile,
				"url":              wrapper.UrlPattern,
			})
		} else {
			_, err = session.Run(ctx, "MATCH (a:App)"+
				"MERGE (b:Servlet {name: $name, class: $class, initParamsKeys: $initParamsKeys, initParamsValues: $initParamsValues, jspFile: $jspFile})"+
				"MERGE (a)-[:route]->(b)", map[string]interface{}{
				"name":             wrapper.ServletName,
				"class":            wrapper.ServletClass,
				"initParamsKeys":   wrapper.InitParamsKeys,
				"initParamsValues": wrapper.InitParamsValues,
				"jspFile":          wrapper.JSPFile,
				"url":              wrapper.UrlPattern,
			})
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Servlet  Name: %s, Class: %s\n", wrapper.ServletName, wrapper.ServletClass)
	}

	return

}
