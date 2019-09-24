package cmd

import (
	"bugbounty/logger"
	"bugbounty/util"
	"crypto/tls"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

var ignoreErrors bool

func init() {
	GetCmd.PersistentFlags().StringVarP(&path, "path", "", "", "path to csv file with domains")
	GetCmd.PersistentFlags().StringVarP(&domains, "domains", "", "", "comma-separated list of domains")
	GetCmd.PersistentFlags().BoolVarP(&ignoreErrors, "ignore-errors", "", true, "ignoreErrors")
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Aliases: []string{"g"},
	Run: func(cmd *cobra.Command, args []string) {
		var hosts []string

		if domains != "" {
			hosts = strings.Split(domains, ",")
		} else if path != "" {
			hosts = util.GetListOfDomainsFromFile(path)
		} else {
			logger.DefaultLogger.Fatal("You need to specify domains or path")
		}

		ch := make(chan string)
		for _, host := range hosts {
			go get(host, ch, ignoreErrors)
		}

		for range hosts {
			resp := <- ch
			if resp != "" {
				logger.DefaultLogger.Info(resp)
			}
		}
	},
}

func get(domain string, out chan<- string, ignoreErrors bool) {
	resp, body, httpErrors := gorequest.
		New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Get(fmt.Sprintf("https://%s", domain)).
		Timeout(3 * time.Second).
		End()

	for _, err := range httpErrors {
		if os.IsTimeout(err) {
			resp, body, httpErrors = gorequest.
				New().
				TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
				Get(fmt.Sprintf("http://%s", domain)).
				Timeout(3 * time.Second).
				End()
		}
	}

	if httpErrors != nil && len(httpErrors) > 0 {
		if ignoreErrors {
			out <- ""
		} else if Verbose {
			out <- fmt.Sprintf("%s: %+v", domain, httpErrors)
		} else {
			out <- fmt.Sprintf("%s: %s", domain, "error")
		}
		return
	}

	var msg string
	if Verbose {
		msg = fmt.Sprintf("%s: %d, body: %s", domain, resp.StatusCode, body)
	} else {
		msg = fmt.Sprintf("%s: %d", domain, resp.StatusCode)
	}

	out <- msg
}