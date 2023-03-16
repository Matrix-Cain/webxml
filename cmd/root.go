package cmd

import (
	"github.com/spf13/cobra"
	"webxml/graph"
)

var (
	neo4jAddr   string
	neo4jUser   string
	neo4jPass   string
	xmlFile     string
	projectName string
	pathVerbose bool

	rootCmd = &cobra.Command{
		Use:   "webxml",
		Short: "A webxml visualizer",
		Long:  `A study project on neo4j`,
		Run: func(cmd *cobra.Command, args []string) {
			graph.BuildGraph(neo4jAddr, neo4jUser, neo4jPass, xmlFile, projectName, pathVerbose)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&neo4jAddr, "addr", "a", "", "neo4j server addr")
	rootCmd.PersistentFlags().StringVarP(&neo4jUser, "user", "u", "", "neo4j server username")
	rootCmd.PersistentFlags().StringVarP(&neo4jPass, "pass", "p", "", "neo4j server password")
	rootCmd.PersistentFlags().StringVarP(&xmlFile, "file", "f", "", "webxml path")
	rootCmd.PersistentFlags().StringVarP(&projectName, "name", "n", "", "project name")
	rootCmd.PersistentFlags().BoolVarP(&pathVerbose, "verbose", "v", false, "verbose path in graph")

}
