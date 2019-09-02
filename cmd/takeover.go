package cmd

import (
	"context"
	"fdns/fdns"
	"fdns/logger"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

func init() {
	Takeover.PersistentFlags().StringVarP(&path, "path", "", "", "path to file")
	Takeover.PersistentFlags().StringVarP(&kind, "kind", "", "", "kind (A|CNAME|NS|PTR)")
}

var Takeover = &cobra.Command{
	Use:     "takeover",
	Short:   "",
	Aliases: []string{"t"},
	Run: func(cmd *cobra.Command, args []string) {
		logger.DefaultLogger.Debug("Starting to scan for takeover")

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

		parser := fdns.NewTakeoverParser(fdns.Domains, 50, parsers[kind])
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

		go parser.ParseTakeover(ctx, r, out, errs)
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
