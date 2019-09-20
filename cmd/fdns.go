package cmd

import (
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	FdnsCmd.AddCommand(Takeover)
	FdnsCmd.AddCommand(Domain)
}

var FdnsCmd = &cobra.Command{
	Use:   "fdns",
	Short: "fdns",
	Long: strings.Join([]string{
		"fdns",
		"",
	}, "\n"),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}
