// Package parser defines interfaces to be implemented by
// parser plugins, used by worker package.
//
package parser

import (
	"fmt"
	"sync"
)

var (
	parsersMu sync.RWMutex
	parsers   = make(map[string]Parser)
)

type Parser interface {
	Parse(url string, target interface{}) error
}

// Register makes a parser driver available by the provided name.
// If Register is called twice with the same name or if parser is nil,
// it panics.
func RegisterParser(name string, parser Parser) {
	parsersMu.Lock()
	defer parsersMu.Unlock()
	if parser == nil {
		panic("Parser: Register parser is nil")
	}
	if _, dup := parsers[name]; dup {
		panic("Parser: Register called twice for parser " + name)
	}
	parsers[name] = parser
}

// Gets registered parser by it's name, if it's not present returns nil
func Get(name string) (Parser, error) {
	parsersMu.RLock()
	parser, ok := parsers[name]
	parsersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("Parser: unknown driver %q (forgotten import?)", name)
	}
	return parser, nil
}
