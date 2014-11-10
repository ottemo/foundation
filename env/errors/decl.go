package errors

import (
	"github.com/ottemo/foundation/env"
	"regexp"
)

const (
	COLLECT_STACK = true
)

var (
	MSG_REGEXP = regexp.MustCompile(`\s*[\[{(]?\s*(?:([0-9]+)?[-: ]([0-9a-fA-F]+)?)?\s*[\]})]?\s*[:\->]*\s*(.+)`)
)

type DefaultErrorBus struct {
	listeners []env.F_ErrorListener
}

type OttemoError struct {
	Message string
	Code    string
	Level   int

	Stack string

	handled bool
}
