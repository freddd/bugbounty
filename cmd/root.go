package cmd

import (
	"fdns/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"strings"
)

var path string
var kind string

func init() {
	RootCmd.PersistentFlags().StringVarP(&path, "path", "", "", "path to file")
	RootCmd.PersistentFlags().StringVarP(&kind, "kind", "", "", "kind (A|CNAME|NS|PTR)")
}


var RootCmd = &cobra.Command{
	Use:   "fdns",
	Short: "fdns",
	Long: strings.Join([]string{
		"fdns",
		"",
	}, "\n"),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(bitbucket completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(bitbucket completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = RootCmd.GenBashCompletionFile("fdns_completions.sh")
	},
}

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generates docs",
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(RootCmd, "./docs")
		if err != nil {
			logger.DefaultLogger.Fatal("%+v", err)
		}
	},
}

func Execute() {
	RootCmd.AddCommand(docsCmd)
	RootCmd.AddCommand(completionCmd)
	RootCmd.AddCommand(Takeover)
	RootCmd.AddCommand(Domain)
	if err := RootCmd.Execute(); err != nil {
		logger.DefaultLogger.Fatal("%s", err.Error())
	}
}
