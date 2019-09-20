package cmd

import (
	"bugbounty/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"strings"
)

var path string
var kind string
var Verbose bool

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

var RootCmd = &cobra.Command{
	Use:   "bugbounty",
	Short: "bb",
	Long: strings.Join([]string{
		"bugbounty",
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
	RootCmd.AddCommand(FdnsCmd)
	RootCmd.AddCommand(Xff)
	RootCmd.AddCommand(ThirdParty)
	if err := RootCmd.Execute(); err != nil {
		logger.DefaultLogger.Fatal("%s", err.Error())
	}
}
