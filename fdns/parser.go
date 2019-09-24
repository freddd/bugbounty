package fdns

import (
	"errors"
	"fmt"
)

// ErrWrongType is the error returned when the parsed entry contains an invalid
// record type.
var ErrWrongType = errors.New("incorrect record type")

// ParseFunc defines how a parsing function must work.
type ParseFunc func(e entry) (string, error)

// A reports DNS A records for the given domain.
func A(e entry) (string, error) {
	if e.Type != "a" {
		return "", ErrWrongType
	}
	return fmt.Sprintf("%s,%s", e.Name, e.Value), nil
}

// CNAME reports DNS CNAME records for the given domain.
func CNAME(e entry) (string, error) {
	if e.Type != "cname" {
		return "", ErrWrongType
	}
	return fmt.Sprintf("%s,%s", e.Name, e.Value), nil
}

// NS reports DNS NS records for the given domain.
func NS(e entry) (string, error) {
	if e.Type != "ns" {
		return "", ErrWrongType
	}
	return e.Name, nil
}

// PTR reports DNS PTR records for the given domain.
func PTR(e entry) (string, error) {
	if e.Type != "ptr" {
		return "", ErrWrongType
	}
	return e.Name, nil
}

// Parser object allows parsing datasets looking for records related with a domain.
type Parser struct {
	domains []string
	// parse defines how the parser looks for results.
	parse ParseFunc
	// workers is the numer of simultaneous goroutines the parser will use.
	workers int

	prefix   string
	suffix   string
	contains string
	regexp   string
	value 	 bool
}
