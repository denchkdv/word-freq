package filedb

import (
	"fmt"
	"hash/fnv"
	"os"
	"path"

	"github.com/denchkdv/word-freq/internal/pkg/database"
	"github.com/denchkdv/word-freq/internal/pkg/wordmap"
)

var _ database.Database = &FileDB{}

type FileDB struct {
	path      string
	volume    uint
	datafiles []*Datafile
}

func New(path string) (*FileDB, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	return &FileDB{path: path}, nil
}

func (db *FileDB) Init(volume uint) error {
	db.volume = volume
	db.datafiles = make([]*Datafile, 0, volume)

	for i := uint(0); i < volume; i++ {
		df, err := NewDatafile(path.Join(db.path, fmt.Sprintf("volume_%d", i)))
		if err != nil {
			return err
		}

		db.datafiles = append(db.datafiles, df)
	}

	return nil
}

func (db *FileDB) Reader(shard uint) database.Reader {
	df := db.datafiles[shard]
	df.MoveStart()
	return df
}

func (db *FileDB) WriteBatch(batch wordmap.WordMap) error {
	blocks := make([][]wordmap.Row, db.volume)

	for key, val := range batch {
		shard := db.shard(key)

		blocks[shard] = append(blocks[shard], wordmap.Row{
			Word:    key,
			Counter: val,
		})
	}

	for i, block := range blocks {
		err := db.datafiles[i].WriteBlock(block)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *FileDB) shard(word string) uint {
	return uint(hash(word)) % db.volume
}

func (db *FileDB) Clean() {
	for _, datafile := range db.datafiles {
		datafile.Remove()
	}
	_ = os.RemoveAll(db.path)
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}
