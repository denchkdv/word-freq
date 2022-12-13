package main

import (
	"flag"
	"fmt"

	"github.com/denchkdv/word-freq/internal/pkg/analyzer"
	"github.com/denchkdv/word-freq/internal/pkg/worker"
)

const (
	databaseDir = "database"
	sampleSize  = analyzer.Sample10Percent
)

func main() {
	input := flag.String("input", "in.txt", "input file")
	memorySize := flag.Int("n", 10, "memory size")
	output := flag.String("output", "out.txt", "output file")

	flag.Parse()

	wrkr, err := worker.New(worker.Params{
		InputPath:   *input,
		OutputPath:  *output,
		DatabaseDir: databaseDir,
		MemorySize:  *memorySize,
		Sample:      sampleSize,
	})
	if err != nil {
		fmt.Printf("Failed to create worker: %v\n", err)
		return
	}

	if err := wrkr.Do(); err != nil {
		fmt.Printf("Failed to process input: %v", err)
	}
}
