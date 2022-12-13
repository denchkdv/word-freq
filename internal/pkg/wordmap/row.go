package wordmap

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type Row struct {
	Word    string
	Counter uint
}

func (r *Row) String() string {
	return fmt.Sprintf("%s\t%d\n", r.Word, r.Counter)
}

var tsvPattern = regexp.MustCompile(`(.+)\t(\d+)`)

func (r *Row) FromString(s string) error {
	submatch := tsvPattern.FindStringSubmatch(s)
	if len(submatch) < 3 {
		return errors.New("unable to parse row")
	}

	counter, err := strconv.ParseUint(submatch[2], 10, 32)
	if err != nil {
		return errors.Wrap(err, "unable to parse word counter")
	}

	r.Word = submatch[1]
	r.Counter = uint(counter)

	return nil
}
