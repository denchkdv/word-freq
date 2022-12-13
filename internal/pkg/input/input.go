package input

import (
	"io"
)

var (
	ErrEOF = io.EOF
)

type Input interface {
	ReadNext() (string, error)
	Size() uint64
	Close() error
}
