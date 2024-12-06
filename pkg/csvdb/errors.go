package csvdb

import (
	"errors"
	"fmt"
)

type ErrorOffset struct {
	actual   int64
	expected int64
}

func (e ErrorOffset) Error() string {
	return fmt.Sprintf("Offset is %d instead of expected %d", e.actual, e.expected)
}

var ErrorFormat = errors.New("CSV format violation")

var ErrorFileTooSmall = errors.New("File is too small to contain data")
