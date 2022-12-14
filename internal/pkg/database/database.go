package database

import "github.com/denchkdv/word-freq/internal/pkg/wordmap"

type Database interface {
	WriteBatch(batch wordmap.WordMap) error
	Reader(shard uint) Reader
}

type Datafile interface {
	Reader
	Writer
}

type Reader interface {
	ReadRow() (wordmap.Row, error)
}

type Writer interface {
	WriteBlock(block []wordmap.Row) error
	Write(wm wordmap.WordMap) error
}
