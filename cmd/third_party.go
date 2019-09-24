package cmd

import (
	"bugbounty/logger"
	"bugbounty/third_party"
	"bugbounty/util"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
)

var target string

func init() {
	ThirdParty.Flags().StringVarP(&domains, "domains", "", "", "comma separated domains")
	ThirdParty.Flags().StringVarP(&path, "path", "", "", "file path")
	ThirdParty.Flags().StringVarP(&target, "target", "", "", "target")
}

var ThirdParty = &cobra.Command{
	Use:     "third-party",
	Short:   "tp",
	Aliases: []string{"tp"},
	Run: func(cmd *cobra.Command, args []string) {
		targets := third_party.Targets
		if target != "" {
			targets = map[string]third_party.Target{
				target: third_party.Targets[target],
			}
		}

		var hosts []string
		if domains != "" {
			hosts = strings.Split(domains, ",")
		} else if path != "" {
			hosts = util.GetListOfDomainsFromFile(path)
		}

		out := make(chan string)
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
					//logger.DefaultLogger.Info("Received SIGINT")
					cancel()
				case <-ctx.Done():
					return
				}
			}
		}()

		go executeScan(hosts, targets, 10, ctx, out)
		go func() {
			for c := range out {
				logger.DefaultLogger.Debug(c)
			}
			done <- struct{}{}
		}()

		<-done
	},
}

func executeScan(hosts []string, targets map[string]third_party.Target, workers int, ctx context.Context, out chan<- string) {
	defer close(out)

	domains := make(chan string)
	done := make(chan struct{})
	finished := make(chan struct{}, workers)

	for i := 0; i < workers; i++ {
		go func() {

			for {
				select {
				case <-done:
					finished <- struct{}{}
					return
				case domain := <-domains:
					for _, target := range targets {
						for _, vulnerability := range target.Vulnerabilities {
							vulnerability.Check(domain)
						}
					}
					out <- fmt.Sprintf("%s: done", domain)
				}
			}
		}()
	}

	var current int
	for _, domain := range hosts {
		select {
		case <-ctx.Done():
			break
		default: // avoid blocking.
		}

		domains <- domain
		current++
	}

	close(done)

	for i := 0; i < workers; i++ {
		<-finished
	}
}
