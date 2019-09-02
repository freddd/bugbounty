package fdns

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)


// Should refactor
func (p *Parser) ParseSubdomains(ctx context.Context, r io.Reader, out chan<- string, errs chan<- error) {
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
						if strings.HasSuffix(e.Name, domain) && !strings.Contains(e.Name, "wpengine") && !strings.Contains(e.Name, "eventbrite") {
							matchesDomain = true
							break
						}
					}

					if !matchesDomain {
						continue
					}

					if p.prefix != "" && !strings.HasPrefix(e.Name, p.prefix) {
						continue
					} else  if p.suffix != "" && !strings.HasSuffix(e.Name, p.suffix) {
						continue
					} else if p.contains != "" && !strings.Contains(e.Name, p.contains) {
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
func NewSubdomainParser(domains []string, workers int, f ParseFunc, prefix string, suffix string, contains string) *Parser {
	return &Parser{domains: domains, parse: f, workers: workers, prefix: prefix, suffix: suffix, contains: contains}
}
