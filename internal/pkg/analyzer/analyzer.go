package analyzer

import (
	"github.com/pkg/errors"

	"github.com/denchkdv/word-freq/internal/pkg/input"
	"github.com/denchkdv/word-freq/internal/pkg/input/file"
)

type Analyzer struct {
	input input.Input
}

func New(path string) (*Analyzer, error) {
	in, err := file.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %#q", path)
	}

	return &Analyzer{
		input: in,
	}, nil
}

type Stats struct {
	Rows uint64
}

type Sample uint64

const (
	Sample1Percent  Sample = 100
	Sample5Percent  Sample = 20
	Sample10Percent Sample = 10
	Sample50Percent Sample = 2
)

func (a *Analyzer) Get(sample Sample) (Stats, error) {
	defer a.input.Close()

	sampleSize := a.calcSampleSize(sample)

	bytesRead, rowsRead := uint64(0), uint64(0)
	for bytesRead < sampleSize {
		s, err := a.input.ReadNext()

		if err == nil {
			rowsRead++
			bytesRead += uint64(len(s)) + 1 // +1 for \n symbol
		} else if errors.Is(err, input.ErrEOF) {
			break
		} else {
			return Stats{}, err
		}
	}

	return Stats{
		Rows: rowsRead * uint64(sample),
	}, nil
}

func (a *Analyzer) calcSampleSize(sampleSize Sample) uint64 {
	inputSize := a.input.Size()

	sample := inputSize / uint64(sampleSize)
	if sample == 0 {
		sample = inputSize
	}

	return sample
}
