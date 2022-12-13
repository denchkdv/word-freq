package worker

import (
	"github.com/pkg/errors"

	"github.com/denchkdv/word-freq/internal/pkg/analyzer"
	"github.com/denchkdv/word-freq/internal/pkg/database"
	"github.com/denchkdv/word-freq/internal/pkg/database/filedb"
	"github.com/denchkdv/word-freq/internal/pkg/input"
	"github.com/denchkdv/word-freq/internal/pkg/input/file"
	"github.com/denchkdv/word-freq/internal/pkg/wordmap"
)

type Params struct {
	InputPath   string
	OutputPath  string
	DatabaseDir string
	MemorySize  int
	Sample      analyzer.Sample
}

type Worker struct {
	params   Params
	input    input.Input
	analyzer *analyzer.Analyzer
	database *filedb.FileDB
	output   database.Writer
}

func New(params Params) (*Worker, error) {
	in, err := file.Open(params.InputPath)
	if err != nil {
		return nil, errors.Wrapf(err, "can't open file %#q", params.InputPath)
	}

	anlzr, err := analyzer.New(params.InputPath)
	if err != nil {
		return nil, errors.Wrap(err, "can't create file analyzer")
	}

	db, err := filedb.New(params.DatabaseDir)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create database in the directory %#q", params.DatabaseDir)
	}

	out, err := filedb.NewDatafile(params.OutputPath)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create output file %#q", params.OutputPath)
	}

	return &Worker{
		params:   params,
		input:    in,
		analyzer: anlzr,
		database: db,
		output:   out,
	}, nil
}

func (w *Worker) estimateDatabaseVolume() (uint, error) {
	stats, err := w.analyzer.Get(w.params.Sample)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to calculate file sample statistics, err: %v")
	}

	// TODO: calc volumes count with respect to word frequency in a sample
	volumes := 1 + uint(stats.Rows/uint64(w.params.MemorySize))

	return volumes, nil
}

func (w *Worker) Do() error {
	volumes, err := w.estimateDatabaseVolume()
	if err != nil {
		return errors.Wrapf(err, "failed to estimate the volume size")
	}

	err = w.database.Init(volumes)
	if err != nil {
		return errors.Wrapf(err, "failed to init the database with %d volumes", volumes)
	}
	defer w.database.Clean()

	err = w.splitInput()
	if err != nil {
		return err
	}

	return w.writeResult(volumes)
}

func (w *Worker) splitInput() error {
	eof := false
	for !eof {
		wm := wordmap.New(w.params.MemorySize)

		for wm.Len() < w.params.MemorySize {
			word, err := w.input.ReadNext()

			if err == nil {
				wm.AddWord(word)
			} else if errors.Is(err, input.ErrEOF) {
				eof = true
				break
			} else {
				return errors.Wrap(err, "failed to read the input file")
			}
		}

		if err := w.database.WriteBatch(wm); err != nil {
			return errors.Wrap(err, "failed to write batch")
		}
	}

	return nil
}

func (w *Worker) writeResult(volumes uint) error {
	for i := uint(0); i < volumes; i++ {
		reader := w.database.Reader(i)
		wm := wordmap.New(w.params.MemorySize)

		eof := false
		// TODO: split the datafile into parts in case of exceeding MemorySize
		for !eof {
			row, err := reader.ReadRow()

			if err == nil {
				wm.AddRow(row)
			} else if errors.Is(err, input.ErrEOF) {
				eof = true
				break
			} else {
				return errors.Wrapf(err, "failed to read the datafile %d", i)
			}
		}

		err := w.output.Write(wm)
		if err != nil {
			return errors.Wrap(err, "failed to write the output")

		}
	}

	return nil
}
