package cmd

import (
	"bugbounty/fdns"
	"bugbounty/logger"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var prefix string
var suffix string
var domain string
var contains string
var regexp string
var value bool

func init() {
	Domain.Flags().StringVarP(&domain, "domain", "", "", "domain")
	Domain.Flags().StringVarP(&prefix, "prefix", "", "", "prefix")
	Domain.Flags().StringVarP(&suffix, "suffix", "", "", "suffix")
	Domain.Flags().StringVarP(&contains, "contains", "", "", "contains")
	Domain.Flags().StringVarP(&regexp, "regexp", "", "", "regexp")
	Domain.Flags().BoolVarP(&value, "value", "", false, "value")
	Domain.Flags().StringVarP(&path, "path", "", "", "path to file")
	Domain.Flags().StringVarP(&kind, "kind", "", "", "kind (A|CNAME|NS|PTR)")
}

var Domain = &cobra.Command{
	Use:     "domain",
	Short:   "",
	Aliases: []string{"d"},
	Run: func(cmd *cobra.Command, args []string) {
		logger.DefaultLogger.Debug("Starting to scan for domain: prefix: %s, suffix: %s, contains: %s, regexp: %s", prefix, suffix, contains, regexp)

		r, err := os.Open(path)
		if err != nil {
			logger.DefaultLogger.Fatal(fmt.Sprintf("could not open file %s: %v", path, err))
		}

		var parsers = map[string]fdns.ParseFunc{
			"A":     fdns.A,
			"CNAME": fdns.CNAME,
			"NS":    fdns.NS,
			"PTR":   fdns.PTR,
		}

		var domains = fdns.Domains
		if domain != "" {
			domains = []string{domain}
		}

		parser := fdns.NewSubdomainParser(domains, 50, parsers[kind], prefix, suffix, contains, regexp, value)
		out := make(chan string)
		errs := make(chan error)
		done := make(chan struct{})

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		c := make(chan os.Signal)
		defer close(c)
		signal.Notify(c, os.Interrupt)

		go func() {
			for {
				select {
				case <-c:
					logger.DefaultLogger.Info("Received SIGINT")
					cancel()
				case <-ctx.Done():
					return
				}
			}
		}()

		go parser.ParseSubdomains(ctx, r, out, errs)
		go func() {
			for {
				select {
				case err := <-errs:
					logger.DefaultLogger.Fatal(fmt.Sprintf("could not parse: %v", err), err)
				}
			}
		}()
		go func() {
			for c := range out {
				logger.DefaultLogger.Info(c)
			}
			done <- struct{}{}
		}()

		<-done
	},
}
