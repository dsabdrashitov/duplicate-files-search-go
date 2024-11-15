package csvdb

import (
	"fmt"
)

type ErrorOffset struct {
	actual   int64
	expected int64
}

func (e ErrorOffset) Error() string {
	return fmt.Sprintf("Offset is %d instead of expected %d", e.actual, e.expected)
}

type ErrorFormat struct {
	cause string
}

func (e ErrorFormat) Error() string {
	return e.cause
}

var ErrorNoNewline = ErrorFormat{"No newline found"}
var ErrorNoQuote = ErrorFormat{"No quote found"}
var ErrorNoComma = ErrorFormat{"No comma found"}
