package filedb

import (
	"bufio"
	"io"
	"os"

	"github.com/denchkdv/word-freq/internal/pkg/database"
	"github.com/denchkdv/word-freq/internal/pkg/input"
	"github.com/denchkdv/word-freq/internal/pkg/wordmap"
)

var _ database.Datafile = &Datafile{}

type Datafile struct {
	path    string
	file    *os.File
	scanner *bufio.Scanner
}

func NewDatafile(path string) (*Datafile, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &Datafile{
		path:    path,
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (df *Datafile) MoveStart() {
	_, _ = df.file.Seek(0, io.SeekStart)
}

func (df *Datafile) ReadRow() (wordmap.Row, error) {
	if !df.scanner.Scan() {
		return wordmap.Row{}, input.ErrEOF
	}

	row := &wordmap.Row{}
	err := row.FromString(df.scanner.Text())

	return *row, err
}

func (df *Datafile) WriteBlock(block []wordmap.Row) error {
	for _, row := range block {
		_, err := df.file.WriteString(row.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (df *Datafile) Write(wm wordmap.WordMap) error {
	for word, counter := range wm {
		_, err := df.file.WriteString((&wordmap.Row{
			Word:    word,
			Counter: counter,
		}).String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (df *Datafile) Remove() {
	_ = os.Remove(df.path)
}
