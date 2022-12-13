package file

import (
	"bufio"
	"os"

	"go.uber.org/multierr"

	"github.com/denchkdv/word-freq/internal/pkg/input"
)

var _ input.Input = &File{}

type File struct {
	info    os.FileInfo
	file    *os.File
	scanner *bufio.Scanner
}

func Open(path string) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &File{
		info:    info,
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (f *File) ReadNext() (string, error) {
	if !f.scanner.Scan() {
		return "", input.ErrEOF
	}

	return f.scanner.Text(), nil
}

func (f *File) Size() uint64 {
	return uint64(f.info.Size())
}

func (f *File) Close() error {
	var err error

	if f.file != nil {
		err = f.file.Close()
	}

	if f.scanner != nil {
		err = multierr.Append(err, f.scanner.Err())
	}

	return err
}
