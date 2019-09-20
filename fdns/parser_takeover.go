package fdns

import (
	"bufio"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Should refactor
func (p *Parser) ParseTakeover(ctx context.Context, r io.Reader, out chan<- string, errs chan<- error) {
	defer close(out)

	gz, err := gzip.NewReader(r)
	if err != nil {
		errs <- err
		return
	}
	defer gz.Close()

	lines := make(chan []byte)
	done := make(chan struct{})
	finished := make(chan struct{}, p.workers)

	for i := 0; i < p.workers; i++ {
		go func() {
			var e entry

			for {
				select {
				case <-done:
					finished <- struct{}{}
					return
				case v := <-lines:
					if err := json.Unmarshal(v, &e); err != nil {
						errs <- fmt.Errorf("could not decode JSON object: %v", err)
						continue
					}

					var matchesDomain = false
					for _, domain := range p.domains {
						if strings.HasSuffix(e.Name, domain) {
							matchesDomain = true
							break
						}
					}

					if !matchesDomain {
						continue
					}

					var s *Site
					for _, site := range VulnerableSites {
						if site.SiteRegex.MatchString(e.Value) {
							s = &site
							break
						}
					}

					if s == nil {
						continue
					}

					rec, err := p.parse(e)
					if err == ErrWrongType {
						// it's not the interesting record type.
						continue
					}
					if err != nil {
						errs <- fmt.Errorf("could not parse object: %v", err)
						continue
					}

					if s.Verify.Kind == "HTTP" {
						resp, body, httpErrors := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).Get(fmt.Sprintf("https://%s", e.Name)).End()

						for _, err := range httpErrors {
							if os.IsTimeout(err) {
								resp, body, httpErrors = gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).Get(fmt.Sprintf("http://%s", e.Name)).End()
							}
						}

						if httpErrors != nil || resp.StatusCode != s.Verify.Condition {
							if resp != nil && resp.StatusCode/100 != 2 && resp.StatusCode != 401 && resp.StatusCode != 403 && resp.StatusCode != 503 && resp.StatusCode != 502 {
								//logger.DefaultLogger.Debug("Name: %s (%s), Status code: %d, body: %s", e.Name, e.Value, resp.StatusCode, body)
							}

							continue
						}

						if !s.Verify.Regexp.MatchString(body) {
							continue
						}
					} else if s.Verify.Kind == "HOST" {
						output, err := exec.Command("sh", "-c", fmt.Sprintf("host %s", e.Name)).CombinedOutput()
						if err == nil {
							continue
						}
						if !s.Verify.Regexp.MatchString(string(output)) {
							continue
						}
					}

					out <- rec
				}
			}
		}()
	}

	sc := bufio.NewScanner(gz)
	var current int
	for sc.Scan() {
		select {
		case <-ctx.Done():
			break
		default: // avoid blocking.
		}

		lines <- append([]byte{}, sc.Bytes()...)
		current++
	}

	if err := sc.Err(); err != nil {
		errs <- fmt.Errorf("could not scan: %v", err)
		return
	}
	close(done)

	for i := 0; i < p.workers; i++ {
		<-finished
	}
}

// NewTakeoverParser returns a FDNS parser that reports entries for the given record.
func NewTakeoverParser(domains []string, workers int, f ParseFunc) *Parser {
	return &Parser{domains: domains, parse: f, workers: workers}
}
