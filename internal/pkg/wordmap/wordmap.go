package wordmap

type WordMap map[string]uint

func New(cap int) WordMap {
	return make(map[string]uint, cap)
}

func (wm WordMap) AddWord(word string) {
	wm[word] = wm[word] + 1
}

func (wm WordMap) AddRow(row Row) {
	wm[row.Word] = wm[row.Word] + row.Counter
}

func (wm WordMap) Len() int {
	return len(wm)
}
